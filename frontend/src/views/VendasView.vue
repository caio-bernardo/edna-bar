<script setup>
import { ref, onMounted, computed, reactive, watch } from "vue";
import api from "@/services/api";

// --- ESTADO GERAL ---
const modo = ref("nova"); // 'nova' ou 'historico'
const carregando = ref(false);

// Caches de dados
const produtos = ref([]);
const clientes = ref([]);
const funcionarios = ref([]);
const lotesCache = ref({}); // Evita buscar o mesmo lote várias vezes

// --- DADOS: NOVA VENDA ---
const carrinho = ref([]);
const formVenda = reactive({
    id_cliente: "",
    id_funcionario: "",
    tipo_pagamento: "pix",
});

// --- DADOS: HISTÓRICO ---
const historicoVendas = ref([]);
const vendaSelecionada = ref(null);
const itensVendaSelecionada = ref([]);

// --- COMPUTED (Cálculos) ---
const totalVenda = computed(() => {
    if (modo.value === "nova") {
        return carrinho.value.reduce(
            (acc, item) => acc + item.preco_venda * item.quantidade,
            0,
        );
    } else if (vendaSelecionada.value) {
        // No histórico, somamos os itens recuperados do banco
        return itensVendaSelecionada.value.reduce(
            (acc, item) => acc + item.valor_unitario * item.quantidade,
            0,
        );
    }
    return 0;
});

// --- AÇÕES GERAIS ---
onMounted(async () => {
    carregarDadosIniciais();
});

const carregarDadosIniciais = async () => {
    try {
        const [resCli, resFunc, resProd] = await Promise.all([
            api.getClientes(),
            api.getFuncionarios(),
            api.getProdutosComerciais(),
        ]);
        clientes.value = resCli.data || [];
        funcionarios.value = resFunc.data || [];
        produtos.value = resProd.data || [];
    } catch (error) {
        console.error("Erro de conexão:", error);
    }
};

// Alternar abas
watch(modo, (novoVal) => {
    if (novoVal === "historico") carregarHistorico();
});

// --- AÇÕES: NOVA VENDA ---
const adicionarAoCarrinho = (produto) => {
    const itemExistente = carrinho.value.find((item) => item.id === produto.id);
    if (itemExistente) {
        itemExistente.quantidade++;
    } else {
        carrinho.value.push({ ...produto, quantidade: 1 });
    }
};

const removerDoCarrinho = (index) => {
    carrinho.value.splice(index, 1);
};

const finalizarVenda = async () => {
    if (!formVenda.id_cliente || !formVenda.id_funcionario)
        return alert("Selecione Cliente e Funcionário.");
    if (carrinho.value.length === 0) return alert("Carrinho vazio.");

    try {
        // 1. Criar Cabeçalho da Venda
        const resVenda = await api.createVenda({
            id_cliente: parseInt(formVenda.id_cliente),
            id_funcionario: parseInt(formVenda.id_funcionario),
            tipo_pagamento: formVenda.tipo_pagamento,
            data_hora_renda: new Date().toISOString(),
            data_hora_pagamento: new Date().toISOString(),
        });

        const idVenda = resVenda.data.id;

        // 2. Processar Itens (Buscar Lote -> Salvar Item)
        for (const item of carrinho.value) {
            // Busca lote disponível (FIFO)
            const resLotes = await api.getLotesPorProduto(item.id);
            const lote = (resLotes.data || [])[0];

            if (lote) {
                await api.createItemVenda({
                    id_venda: idVenda,

                    id_lote: lote.id_lote,
                    quantidade: parseInt(item.quantidade),
                    valor_unitario: parseFloat(item.preco_venda),
                });
            } else {
                console.warn(`Sem lote para produto ${item.nome}`);
            }
        }

        alert("Venda realizada com sucesso!");
        carrinho.value = [];
        formVenda.id_cliente = "";
    } catch (error) {
        console.error(error);
        alert(
            "Erro ao vender: " +
                (error.response?.data?.detail ||
                    "Verifique a conexão com o backend"),
        );
    }
};

// --- AÇÕES: HISTÓRICO ---
const carregarHistorico = async () => {
    carregando.value = true;
    try {
        // Ordena por data decrescente (backend sort)
        const res = await api.getVendas();
        historicoVendas.value = res.data || [];
    } finally {
        carregando.value = false;
    }
};

const selecionarVendaHistorico = async (venda) => {
    vendaSelecionada.value = venda;
    itensVendaSelecionada.value = [];

    try {
        // 1. Buscar itens da venda
        const resItens = await api.getItemVenda({
            params: { "filter-id_venda": `eq.${venda.id}` },
        });
        const itensRaw = resItens.data || [];

        // 2. Enriquecer com nomes dos produtos (vêm apenas com ID do lote/produto)
        const itensDetalhados = await Promise.all(
            itensRaw.map(async (item) => {
                let nomeProd = "...";

                // Tenta cache do lote
                if (!lotesCache.value[item.id_lote]) {
                    try {
                        const l = await api.getLote(item.id_lote);
                        lotesCache.value[item.id_lote] = l.data;
                    } catch {
                        /* lote deletado? */
                    }
                }

                const lote = lotesCache.value[item.id_lote];
                if (lote) {
                    const prod = produtos.value.find(
                        (p) => p.id === lote.id_produto,
                    );
                    if (prod) nomeProd = prod.nome;
                }

                return { ...item, nome: nomeProd };
            }),
        );

        itensVendaSelecionada.value = itensDetalhados;
    } catch (error) {
        console.error(error);
    }
};

// Helpers visuais
const getNomeCliente = (id) =>
    clientes.value.find((c) => c.id === id)?.nome || "N/A";
const getNomeFunc = (id) =>
    funcionarios.value.find((f) => f.id === id)?.nome || "N/A";
const formatarData = (isoStr) => new Date(isoStr).toLocaleString("pt-BR");
</script>

<template>
    <div class="nav-space"></div>

    <div class="tabs">
        <button :class="{ active: modo === 'nova' }" @click="modo = 'nova'">
            Nova Venda
        </button>
        <button
            :class="{ active: modo === 'historico' }"
            @click="modo = 'historico'"
        >
            Histórico
        </button>
    </div>

    <div class="vendas-layout">
        <div class="panel-left">
            <div
                class="receipt-paper"
                :class="{ 'receipt-readonly': modo === 'historico' }"
            >
                <h2 class="receipt-title">
                    {{
                        modo === "nova"
                            ? "Caixa Aberto"
                            : vendaSelecionada
                              ? `Nota #${vendaSelecionada.id}`
                              : "Selecione"
                    }}
                </h2>

                <div class="receipt-header">
                    <div v-if="modo === 'nova'" class="form-stack">
                        <div class="form-group">
                            <label>Cliente</label>
                            <select v-model="formVenda.id_cliente">
                                <option value="" disabled>Selecione...</option>
                                <option
                                    v-for="c in clientes"
                                    :key="c.id"
                                    :value="c.id"
                                >
                                    {{ c.nome }}
                                </option>
                            </select>
                        </div>
                        <div class="row">
                            <div class="form-group grow">
                                <label>Atendente</label>
                                <select v-model="formVenda.id_funcionario">
                                    <option value="" disabled>...</option>
                                    <option
                                        v-for="f in funcionarios"
                                        :key="f.id"
                                        :value="f.id"
                                    >
                                        {{ f.nome }}
                                    </option>
                                </select>
                            </div>
                            <div class="form-group">
                                <label>Pagamento</label>
                                <select v-model="formVenda.tipo_pagamento">
                                    <option value="pix">Pix</option>
                                    <option value="credito">Crédito</option>
                                    <option value="debito">Débito</option>
                                    <option value="dinheiro">Dinheiro</option>
                                </select>
                            </div>
                        </div>
                    </div>

                    <div v-else-if="vendaSelecionada" class="info-static">
                        <p>
                            <span>Cliente:</span>
                            {{ getNomeCliente(vendaSelecionada.id_cliente) }}
                        </p>
                        <p>
                            <span>Atendente:</span>
                            {{ getNomeFunc(vendaSelecionada.id_funcionario) }}
                        </p>
                        <p>
                            <span>Data:</span>
                            {{ formatarData(vendaSelecionada.data_hora_renda) }}
                        </p>
                        <p>
                            <span>Pgto:</span>
                            {{ vendaSelecionada.tipo_pagamento.toUpperCase() }}
                        </p>
                    </div>
                    <div v-else class="msg-empty">
                        Selecione uma venda da lista &rarr;
                    </div>
                </div>

                <div class="receipt-items">
                    <div class="items-head">
                        <span>Qtd</span>
                        <span>Item</span>
                        <span>Valor</span>
                    </div>
                    <div class="items-body">
                        <template v-if="modo === 'nova'">
                            <div
                                v-for="(item, idx) in carrinho"
                                :key="idx"
                                class="item-row"
                                @dblclick="removerDoCarrinho(idx)"
                            >
                                <span class="col-qtd"
                                    >{{ item.quantidade }}x</span
                                >
                                <span class="col-name">{{ item.nome }}</span>
                                <span class="col-price">{{
                                    (
                                        item.preco_venda * item.quantidade
                                    ).toFixed(2)
                                }}</span>
                            </div>
                        </template>
                        <template v-else-if="vendaSelecionada">
                            <div
                                v-for="(item, idx) in itensVendaSelecionada"
                                :key="idx"
                                class="item-row"
                            >
                                <span class="col-qtd"
                                    >{{ item.quantidade }}x</span
                                >
                                <span class="col-name">{{ item.nome }}</span>
                                <span class="col-price">{{
                                    (
                                        item.valor_unitario * item.quantidade
                                    ).toFixed(2)
                                }}</span>
                            </div>
                        </template>
                    </div>
                </div>

                <div class="receipt-total">
                    <span>TOTAL</span>
                    <span class="big-price"
                        >R$ {{ totalVenda.toFixed(2) }}</span
                    >
                </div>

                <button
                    v-if="modo === 'nova'"
                    class="btn-fab"
                    @click="finalizarVenda"
                >
                    +
                </button>
            </div>
        </div>

        <div class="panel-right">
            <div v-if="modo === 'nova'" class="catalog-grid">
                <div
                    v-for="prod in produtos"
                    :key="prod.id"
                    class="card-prod"
                    @click="adicionarAoCarrinho(prod)"
                >
                    <div class="prod-info">
                        <span class="prod-name">{{ prod.nome }}</span>
                        <span class="prod-brand">{{ prod.marca }}</span>
                    </div>
                    <span class="prod-price">R$ {{ prod.preco_venda }}</span>
                </div>
            </div>

            <div v-else class="history-list">
                <div
                    v-for="v in historicoVendas"
                    :key="v.id"
                    class="card-sale"
                    :class="{ active: vendaSelecionada?.id === v.id }"
                    @click="selecionarVendaHistorico(v)"
                >
                    <div class="sale-left">
                        <span class="sale-id">#{{ v.id }}</span>
                        <span class="sale-client">{{
                            getNomeCliente(v.id_cliente)
                        }}</span>
                    </div>
                    <div class="sale-right">
                        <span class="sale-date">{{
                            formatarData(v.data_hora_renda).split(",")[0]
                        }}</span>
                        <span class="sale-tag">{{ v.tipo_pagamento }}</span>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* Layout Base */
.vendas-layout {
    display: flex;
    height: 85vh;
    background-color: var(--edna-black);
    overflow: hidden;
}

/* Abas */
.tabs {
    background-color: var(--edna-black);
    padding: 0 20px;
    border-bottom: 1px solid var(--edna-gray);
    display: flex;
    gap: 10px;
}
.tabs button {
    background: none;
    border: none;
    color: var(--edna-light-gray);
    padding: 15px 20px;
    font-size: 1rem;
    cursor: pointer;
    border-bottom: 3px solid transparent;
}
.tabs button.active {
    color: var(--edna-white);
    border-bottom-color: var(--edna-yellow);
    font-weight: bold;
}

/* Esquerda - Nota Fiscal */
.panel-left {
    width: 350px;
    min-width: 350px;
    padding: 20px;
    background-color: var(--edna-dark-gray);
    border-right: 2px solid var(--edna-wine);
}

.receipt-paper {
    background-color: #1e1e24;
    border: 1px solid var(--edna-yellow);
    border-radius: 6px;
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 15px;
    position: relative;
    box-shadow: 0 0 15px rgba(0, 0, 0, 0.5);
}

.receipt-readonly {
    border-color: var(--edna-light-gray);
    border-style: dashed;
}

.receipt-title {
    text-align: center;
    color: var(--edna-yellow);
    font-family: "IM Fell English SC", serif;
    margin-bottom: 15px;
    border-bottom: 1px double var(--edna-gray);
    padding-bottom: 10px;
}

/* Formulários Nota */
.form-stack {
    display: flex;
    flex-direction: column;
    gap: 10px;
}
.row {
    display: flex;
    gap: 10px;
}
.grow {
    flex: 1;
}
.form-group label {
    display: block;
    font-size: 0.75rem;
    color: var(--edna-light-gray);
    margin-bottom: 3px;
}
select {
    width: 100%;
    background-color: var(--edna-black);
    color: white;
    border: 1px solid var(--edna-gray);
    padding: 6px;
    border-radius: 4px;
}

/* Info Estática */
.info-static p {
    display: flex;
    justify-content: space-between;
    font-size: 0.9rem;
    border-bottom: 1px dotted #333;
    padding: 4px 0;
}
.info-static span {
    color: var(--edna-light-gray);
}
.msg-empty {
    text-align: center;
    color: #555;
    margin-top: 50px;
    font-style: italic;
}

/* Lista Itens Nota */
.receipt-items {
    flex: 1;
    margin: 15px 0;
    border-top: 1px dashed var(--edna-gray);
    border-bottom: 1px dashed var(--edna-gray);
    display: flex;
    flex-direction: column;
}
.items-head {
    display: flex;
    font-size: 0.8rem;
    color: var(--edna-yellow);
    padding: 5px 0;
    font-weight: bold;
}
.items-body {
    flex: 1;
    overflow-y: auto;
}
.item-row {
    display: flex;
    font-size: 0.9rem;
    padding: 4px 0;
    font-family: monospace;
}
.col-qtd {
    width: 40px;
    text-align: center;
    color: var(--edna-light-gray);
}
.col-name {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.col-price {
    width: 70px;
    text-align: right;
    color: var(--edna-white);
}

/* Total e Botão */
.receipt-total {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 1.2rem;
    margin-bottom: 30px;
}
.big-price {
    color: var(--edna-green);
    font-weight: bold;
    font-size: 1.5rem;
}

.btn-fab {
    position: absolute;
    bottom: 15px;
    left: 50%;
    transform: translateX(-50%);
    width: 50px;
    height: 50px;
    border-radius: 50%;
    background-color: var(--edna-yellow);
    color: var(--edna-black);
    font-size: 2rem;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.5);
    transition: transform 0.2s;
}
.btn-fab:hover {
    transform: translateX(-50%) scale(1.1);
    background-color: white;
}

/* Direita */
.panel-right {
    flex: 1;
    padding: 20px;
    overflow-y: auto;
}

/* Grid Catálogo */
.catalog-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 15px;
}
.card-prod {
    background-color: var(--edna-gray);
    border: 1px solid transparent;
    padding: 15px;
    border-radius: 8px;
    cursor: pointer;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    height: 120px;
    transition: all 0.2s;
}
.card-prod:hover {
    border-color: var(--edna-yellow);
    background-color: #3e3e50;
    transform: translateY(-3px);
}
.prod-name {
    font-weight: bold;
    display: block;
}
.prod-brand {
    font-size: 0.8rem;
    color: var(--edna-light-gray);
}
.prod-price {
    align-self: flex-end;
    color: var(--edna-green);
    font-weight: bold;
}

/* Lista Histórico */
.history-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
}
.card-sale {
    background-color: var(--edna-gray);
    padding: 12px;
    border-radius: 6px;
    display: flex;
    justify-content: space-between;
    cursor: pointer;
    border-left: 4px solid transparent;
}
.card-sale:hover {
    background-color: #3e3e50;
}
.card-sale.active {
    border-left-color: var(--edna-yellow);
    background-color: #2c2c35;
}
.sale-id {
    color: var(--edna-yellow);
    font-weight: bold;
    margin-right: 10px;
}
.sale-tag {
    background-color: var(--edna-black);
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 0.7rem;
    text-transform: uppercase;
}
.sale-date {
    color: var(--edna-light-gray);
    font-size: 0.8rem;
    margin-right: 10px;
}
</style>
