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
:root {
    --edna-blue: #b6e5f3;
    --edna-green: #5ad3b0;
    --edna-wine: #a12d4c;
    --edna-red: #e71d51;
    --edna-orange: #f4716e;
    --edna-yellow: #ffd782;
    --edna-light-gray: #888899;
    --edna-gray: #353545;
    --edna-dark-gray: #2a2a32;
    --edna-black: #1a1a1e;
    --edna-white: #f4f4ff;
}

/* --- LAYOUT GERAL --- */
.nav-space {
    background-image: linear-gradient(
        220deg,
        var(--edna-yellow),
        var(--edna-green)
    );
}

.vendas-layout {
    display: flex;
    height: calc(100vh - 16vh); /* Ajuste para não estourar a tela com a nav e tabs */
    background-color: var(--edna-black);
    overflow: hidden;
    font-family: "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
    color: var(--edna-white);
}

/* --- ABAS DE NAVEGAÇÃO --- */
.tabs {
    background-color: var(--edna-dark-gray);
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
    transition: all 0.3s;
}

.tabs button:hover {
    color: var(--edna-white);
    background-color: rgba(255, 255, 255, 0.05);
}

.tabs button.active {
    color: var(--edna-yellow);
    border-bottom-color: var(--edna-yellow);
    font-weight: bold;
}

/* --- COLUNA ESQUERDA (NOTA FISCAL) --- */
.panel-left {
    width: 380px;
    min-width: 350px;
    padding: 20px;
    background-color: var(--edna-dark-gray);
    border-right: 1px solid var(--edna-gray);
    display: flex;
    flex-direction: column;
}

/* Estilo do Papel/Recibo */
.receipt-paper {
    background-color: var(--edna-gray); /* Fundo do card */
    border: 1px solid var(--edna-yellow); /* Borda amarela característica */
    border-radius: 12px;
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 20px;
    position: relative;
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3);
}

.receipt-readonly {
    border-color: var(--edna-light-gray);
    border-style: dashed;
    opacity: 0.9;
}

.receipt-title {
    text-align: center;
    color: var(--edna-yellow);
    font-weight: 300;
    font-size: 1.5rem;
    margin-top: 0;
    margin-bottom: 20px;
    border-bottom: 1px dashed var(--edna-light-gray);
    padding-bottom: 15px;
    text-transform: uppercase;
    letter-spacing: 2px;
}

/* Formulários dentro da Nota */
.form-stack {
    display: flex;
    flex-direction: column;
    gap: 15px;
}

.row {
    display: flex;
    gap: 10px;
}

.grow { flex: 1; }

.form-group label {
    display: block;
    font-size: 0.8rem;
    color: var(--edna-light-gray);
    margin-bottom: 5px;
}

select {
    width: 100%;
    background-color: var(--edna-black);
    color: var(--edna-white);
    border: 1px solid var(--edna-gray);
    padding: 10px;
    border-radius: 6px;
    outline: none;
    box-sizing: border-box; /* Importante para não vazar */
}

select:focus {
    border-color: var(--edna-yellow);
}

/* Info Estática (Histórico) */
.info-static p {
    display: flex;
    justify-content: space-between;
    font-size: 0.95rem;
    border-bottom: 1px dotted var(--edna-light-gray);
    padding: 8px 0;
    margin: 0;
}
.info-static span {
    color: var(--edna-light-gray);
}

.msg-empty {
    text-align: center;
    color: var(--edna-light-gray);
    margin-top: 50px;
    font-style: italic;
}

/* Lista de Itens na Nota */
.receipt-items {
    flex: 1;
    margin: 20px 0;
    border-top: 2px solid var(--edna-black);
    border-bottom: 2px solid var(--edna-black);
    background-color: rgba(0, 0, 0, 0.2);
    border-radius: 4px;
    padding: 10px;
    display: flex;
    flex-direction: column;
    overflow: hidden; /* Contém o scroll */
}

.items-head {
    display: flex;
    font-size: 0.85rem;
    color: var(--edna-yellow);
    padding-bottom: 8px;
    border-bottom: 1px solid var(--edna-gray);
    font-weight: bold;
    text-transform: uppercase;
}

.items-body {
    flex: 1;
    overflow-y: auto;
    padding-top: 10px;
    /* Scrollbar fina */
    scrollbar-width: thin;
    scrollbar-color: var(--edna-gray) var(--edna-black);
}

.item-row {
    display: flex;
    font-size: 0.95rem;
    padding: 6px 0;
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.item-row:hover {
    background-color: rgba(255, 255, 255, 0.05);
}

.col-qtd {
    width: 40px;
    text-align: center;
    color: var(--edna-blue);
    font-weight: bold;
}

.col-name {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    padding: 0 10px;
}

.col-price {
    width: 80px;
    text-align: right;
    color: var(--edna-green);
}

/* Totalizador */
.receipt-total {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: auto;
    padding-top: 15px;
}

.receipt-total span:first-child {
    font-size: 1.2rem;
    color: var(--edna-light-gray);
}

.big-price {
    color: var(--edna-green);
    font-weight: bold;
    font-size: 2rem;
    text-shadow: 0 0 5px rgba(90, 211, 176, 0.2);
}

/* Botão Flutuante (+) */
.btn-fab {
    position: absolute;
    bottom: -25px; /* Metade para fora do card para estilo */
    left: 50%;
    transform: translateX(-50%);
    width: 60px;
    height: 60px;
    border-radius: 50%;
    background-color: var(--edna-green);
    color: var(--edna-black);
    font-size: 2.5rem;
    border: 4px solid var(--edna-dark-gray); /* Borda combina com fundo do painel */
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.5);
    transition: transform 0.2s, background-color 0.2s;
    padding-bottom: 6px; /* Ajuste visual do + */
}

.btn-fab:hover {
    transform: translateX(-50%) scale(1.1);
    background-color: var(--edna-white);
}

/* --- COLUNA DIREITA (CATÁLOGO/HISTÓRICO) --- */
.panel-right {
    flex: 1;
    padding: 20px;
    background-color: var(--edna-black);
    overflow-y: auto;
}

/* Grid de Produtos */
.catalog-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: 20px;
}

.card-prod {
    background-color: var(--edna-dark-gray);
    border: 1px solid var(--edna-gray);
    padding: 15px;
    border-radius: 8px;
    cursor: pointer;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    height: 140px;
    transition: all 0.2s;
    border-left: 4px solid var(--edna-blue); /* Identidade visual de produto */
}

.card-prod:hover {
    background-color: var(--edna-gray);
    transform: translateY(-5px);
    box-shadow: 0 5px 15px rgba(0,0,0,0.3);
}

.prod-info {
    display: flex;
    flex-direction: column;
}

.prod-name {
    font-weight: bold;
    font-size: 1.1rem;
    color: var(--edna-white);
    margin-bottom: 5px;
}

.prod-brand {
    font-size: 0.85rem;
    color: var(--edna-light-gray);
    text-transform: uppercase;
    background-color: var(--edna-black);
    padding: 2px 6px;
    border-radius: 4px;
    align-self: flex-start;
}

.prod-price {
    align-self: flex-end;
    color: var(--edna-green);
    font-weight: bold;
    font-size: 1.2rem;
}

/* Lista Histórico */
.history-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.card-sale {
    background-color: var(--edna-dark-gray);
    padding: 15px 20px;
    border-radius: 8px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    border: 1px solid var(--edna-gray);
    border-left: 4px solid var(--edna-light-gray);
    transition: all 0.2s;
}

.card-sale:hover {
    background-color: var(--edna-gray);
}

.card-sale.active {
    border-left-color: var(--edna-yellow);
    background-color: #3d3d4d;
    border-color: var(--edna-yellow);
}

.sale-left {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.sale-id {
    color: var(--edna-yellow);
    font-weight: bold;
    font-size: 0.9rem;
}

.sale-client {
    font-size: 1.1rem;
    font-weight: bold;
    color: var(--edna-white);
}

.sale-right {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 4px;
}

.sale-date {
    color: var(--edna-light-gray);
    font-size: 0.85rem;
}

.sale-tag {
    background-color: var(--edna-black);
    color: var(--edna-blue);
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 0.8rem;
    text-transform: uppercase;
    font-weight: bold;
}

/* Scrollbars globais para o componente */
::-webkit-scrollbar { width: 8px; }
::-webkit-scrollbar-track { background: var(--edna-black); }
::-webkit-scrollbar-thumb { background: var(--edna-gray); border-radius: 4px; }
::-webkit-scrollbar-thumb:hover { background: var(--edna-light-gray); }
</style>
