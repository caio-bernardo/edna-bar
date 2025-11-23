<script setup>
import { ref, onMounted, computed, reactive, watch } from "vue";
import api from "@/services/api";

// --- ESTADO GERAL ---
const modo = ref("nova"); // 'nova', 'historico', 'clientes'
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

// --- DADOS: CLIENTES ---
const clienteSelecionado = ref(null);
const historicoCliente = ref([]);
const saldoCliente = ref(0);

// --- DADOS: NOVO CLIENTE ---
const isCreatingClient = ref(false);
const formCliente = reactive({
    nome: "",
    cpf: "",
    data_nascimento: "",
});

// --- COMPUTED (Cálculos) ---
const totalVenda = computed(() => {
    if (modo.value === "nova") {
        return carrinho.value.reduce(
            (acc, item) => acc + item.preco_venda * item.quantidade,
            0,
        );
    } else if (modo.value === "historico" && vendaSelecionada.value) {
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

// Alternar abas e resetar estados
watch(modo, (novoVal) => {
    // RESETAR O MODO DE CRIAÇÃO AO TROCAR DE ABA
    isCreatingClient.value = false;

    if (novoVal === "historico") {
        carregarHistorico();
        clienteSelecionado.value = null;
    } else if (novoVal === "clientes") {
        vendaSelecionada.value = null;
        clienteSelecionado.value = null;
    } else if (novoVal === "nova") {
        vendaSelecionada.value = null;
        clienteSelecionado.value = null;
    }
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
        // Ordena por data decrescente (backend sort - snake_case)
        const res = await api.getVendas({
            params: { sort: "-data_hora_venda" },
        });
        historicoVendas.value = res.data || [];
    } catch (error) {
        console.error("Erro ao carregar histórico:", error);
    } finally {
        carregando.value = false;
    }
};

const selecionarVendaHistorico = async (venda) => {
    vendaSelecionada.value = venda;
    itensVendaSelecionada.value = [];

    try {
        const resItens = await api.getItemVenda({
            params: { "filter-id_venda": `eq.${venda.id}` },
        });
        const itensRaw = resItens.data || [];

        // Enriquecer com nomes dos produtos
        const itensDetalhados = await Promise.all(
            itensRaw.map(async (item) => {
                let nomeProd = "...";
                if (!lotesCache.value[item.id_lote]) {
                    try {
                        const l = await api.getLote(item.id_lote);
                        lotesCache.value[item.id_lote] = l.data;
                    } catch {}
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

// --- AÇÕES: CLIENTES ---
const iniciarCadastroCliente = () => {
    clienteSelecionado.value = null;
    isCreatingClient.value = true;
    Object.assign(formCliente, { nome: "", cpf: "", data_nascimento: "" });
};

const salvarCliente = async () => {
    if (!formCliente.nome || !formCliente.cpf)
        return alert("Nome e CPF obrigatórios");

    try {
        const payload = { ...formCliente };
        if (payload.data_nascimento) {
            payload.data_nascimento = new Date(
                payload.data_nascimento,
            ).toISOString();
        } else {
            payload.data_nascimento = null;
        }

        await api.createCliente(payload);

        alert("Cliente cadastrado!");
        isCreatingClient.value = false;
        carregarDadosIniciais(); // Recarrega a lista
    } catch (error) {
        alert(
            "Erro ao criar: " + (error.response?.data?.detail || error.message),
        );
    }
};

const selecionarCliente = async (cliente) => {
    isCreatingClient.value = false;
    clienteSelecionado.value = cliente;
    saldoCliente.value = 0;
    historicoCliente.value = [];

    try {
        // 1. Buscar Saldo (Garante float para exibição)
        const resSaldo = await api.getClienteSaldo(cliente.id);
        saldoCliente.value = parseFloat(resSaldo.data.saldo_devedor || 0);

        // 2. Buscar Histórico
        const resVendas = await api.getVendas({
            params: {
                "filter-id_cliente": `eq.${cliente.id}`,
                sort: "-data_hora_venda",
            },
        });

        const vendasRaw = resVendas.data || [];

        // 3. Calcular TOTAL de cada venda para o extrato
        const historicoComValores = await Promise.all(
            vendasRaw.map(async (v) => {
                // Busca itens da venda para somar
                const resItens = await api.getItemVenda({
                    params: { "filter-id_venda": `eq.${v.id}` },
                });
                const itens = resItens.data || [];

                // Soma: Quantidade * Valor Unitário
                const totalVenda = itens.reduce(
                    (acc, i) => acc + i.quantidade * i.valor_unitario,
                    0,
                );

                return {
                    ...v,
                    total: totalVenda,
                };
            }),
        );

        historicoCliente.value = historicoComValores;
    } catch (error) {
        console.error("Erro ao carregar dados do cliente:", error);
    }
};

// Helpers visuais
const getNomeCliente = (id) =>
    clientes.value.find((c) => c.id === id)?.nome || "N/A";
const getNomeFunc = (id) =>
    funcionarios.value.find((f) => f.id === id)?.nome || "N/A";
const formatarData = (isoStr) =>
    isoStr ? new Date(isoStr).toLocaleString("pt-BR") : "-";
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
        <button
            :class="{ active: modo === 'clientes' }"
            @click="modo = 'clientes'"
        >
            Clientes
        </button>
    </div>

    <div class="vendas-layout">
        <div class="panel-left">
            <div
                class="receipt-paper"
                :class="{ 'receipt-readonly': modo !== 'nova' }"
            >
                <h2 class="receipt-title">
                    <span v-if="modo === 'nova'">Caixa Aberto</span>
                    <span v-else-if="modo === 'historico'">
                        {{
                            vendaSelecionada
                                ? `Nota #${vendaSelecionada.id}`
                                : "Selecione"
                        }}
                    </span>
                    <span v-else-if="modo === 'clientes'">
                        {{
                            isCreatingClient
                                ? "Novo Cadastro"
                                : "Extrato Cliente"
                        }}
                    </span>
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

                    <div v-else-if="modo === 'historico'">
                        <div v-if="vendaSelecionada" class="info-static">
                            <p>
                                <span>Cliente:</span>
                                {{
                                    getNomeCliente(vendaSelecionada.id_cliente)
                                }}
                            </p>
                            <p>
                                <span>Atendente:</span>
                                {{
                                    getNomeFunc(vendaSelecionada.id_funcionario)
                                }}
                            </p>
                            <p>
                                <span>Data:</span>
                                {{
                                    formatarData(
                                        vendaSelecionada.data_hora_renda,
                                    )
                                }}
                            </p>
                            <p>
                                <span>Pgto:</span>
                                {{
                                    vendaSelecionada.tipo_pagamento.toUpperCase()
                                }}
                            </p>
                        </div>
                        <div v-else class="msg-empty">
                            Selecione uma venda da lista &rarr;
                        </div>
                    </div>

                    <div v-else-if="modo === 'clientes'">
                        <div v-if="isCreatingClient" class="form-stack">
                            <div class="form-group">
                                <label>Nome Completo</label>
                                <input
                                    v-model="formCliente.nome"
                                    type="text"
                                    placeholder="Ex: João da Silva"
                                    class="input-dark"
                                />
                            </div>
                            <div class="form-group">
                                <label>CPF (Apenas números)</label>
                                <input
                                    v-model="formCliente.cpf"
                                    type="text"
                                    maxlength="11"
                                    placeholder="00011122233"
                                    class="input-dark"
                                />
                            </div>
                            <div class="form-group">
                                <label>Data de Nascimento</label>
                                <input
                                    v-model="formCliente.data_nascimento"
                                    type="date"
                                    class="input-dark"
                                />
                            </div>
                        </div>

                        <div v-else-if="clienteSelecionado" class="info-static">
                            <p>
                                <span>Nome:</span> {{ clienteSelecionado.nome }}
                            </p>
                            <p>
                                <span>CPF:</span>
                                {{ clienteSelecionado.cpf || "N/A" }}
                            </p>
                            <p>
                                <span>Nascimento:</span>
                                {{
                                    formatarData(
                                        clienteSelecionado.data_nascimento,
                                    ).split(",")[0]
                                }}
                            </p>
                        </div>
                        <div v-else class="msg-empty">
                            Selecione um cliente ou crie um novo &rarr;
                        </div>
                    </div>
                </div>

                <div class="receipt-items" v-if="!isCreatingClient">
                    <div class="items-head" v-if="modo !== 'clientes'">
                        <span>Qtd</span>
                        <span>Item</span>
                        <span>Valor</span>
                    </div>
                    <div class="items-head" v-else>
                        <span>Data</span>
                        <span>Pgto</span>
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

                        <template
                            v-else-if="modo === 'historico' && vendaSelecionada"
                        >
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

                        <template
                            v-else-if="
                                modo === 'clientes' && clienteSelecionado
                            "
                        >
                            <div
                                v-for="v in historicoCliente"
                                :key="v.id"
                                class="item-row"
                            >
                                <span class="col-date">{{
                                    formatarData(v.data_hora_renda).split(
                                        " ",
                                    )[0]
                                }}</span>
                                <span class="col-name">{{
                                    v.tipo_pagamento
                                }}</span>
                                <span class="col-price"
                                    >R$ {{ v.total.toFixed(2) }}</span
                                >
                            </div>
                            <div
                                v-if="historicoCliente.length === 0"
                                class="msg-empty-small"
                            >
                                Sem histórico de compras
                            </div>
                        </template>
                    </div>
                </div>

                <div v-else style="flex: 1"></div>

                <div class="receipt-total" v-if="!isCreatingClient">
                    <template v-if="modo !== 'clientes'">
                        <span>TOTAL</span>
                        <span class="big-price"
                            >R$ {{ totalVenda.toFixed(2) }}</span
                        >
                    </template>
                    <template v-else-if="modo === 'clientes'">
                        <span>SALDO DEVEDOR</span>
                        <span
                            class="big-price"
                            :class="{ debt: saldoCliente > 0 }"
                        >
                            R$ {{ saldoCliente.toFixed(2) }}
                        </span>
                    </template>
                </div>

                <button
                    v-if="modo === 'nova'"
                    class="btn-fab"
                    @click="finalizarVenda"
                >
                    +
                </button>
                <button
                    v-if="modo === 'clientes' && isCreatingClient"
                    class="btn-fab"
                    @click="salvarCliente"
                >
                    ✓
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

            <div v-else-if="modo === 'historico'" class="history-list">
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

            <div v-else-if="modo === 'clientes'" class="history-list">
                <div
                    class="card-sale new-client-card"
                    @click="iniciarCadastroCliente"
                    :class="{ active: isCreatingClient }"
                >
                    <div
                        class="sale-left"
                        style="
                            flex-direction: row;
                            align-items: center;
                            gap: 10px;
                        "
                    >
                        <span class="sale-id" style="font-size: 1.5rem">+</span>
                        <span class="sale-client">Cadastrar Novo Cliente</span>
                    </div>
                </div>

                <div
                    v-for="c in clientes"
                    :key="c.id"
                    class="card-sale card-cliente"
                    :class="{ active: clienteSelecionado?.id === c.id }"
                    @click="selecionarCliente(c)"
                >
                    <div class="sale-left">
                        <span class="sale-client">{{ c.nome }}</span>
                        <span class="client-cpf">{{ c.cpf }}</span>
                    </div>
                    <div class="sale-right">
                        <span class="sale-tag">ID: {{ c.id }}</span>
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
    height: calc(100vh - 16vh);
    background-color: var(--edna-black);
    overflow: hidden;
    font-family: "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
    color: var(--edna-white);
}

/* --- ABAS --- */
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

/* --- PAINEL ESQUERDO --- */
.panel-left {
    width: 380px;
    min-width: 350px;
    padding: 20px;
    background-color: var(--edna-dark-gray);
    border-right: 1px solid var(--edna-gray);
    display: flex;
    flex-direction: column;
}

.receipt-paper {
    background-color: var(--edna-gray);
    border: 1px solid var(--edna-yellow);
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

.form-stack {
    display: flex;
    flex-direction: column;
    gap: 15px;
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
    box-sizing: border-box;
}
select:focus {
    border-color: var(--edna-yellow);
}

.input-dark {
    width: 100%;
    background-color: var(--edna-black);
    color: var(--edna-white);
    border: 1px solid var(--edna-gray);
    padding: 10px;
    border-radius: 6px;
    outline: none;
    box-sizing: border-box;
    font-family: "Segoe UI", sans-serif;
}
.input-dark:focus {
    border-color: var(--edna-yellow);
}

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
.msg-empty-small {
    text-align: center;
    color: var(--edna-light-gray);
    margin-top: 10px;
    font-size: 0.8rem;
}

/* Tabela */
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
    overflow: hidden;
}

.items-head {
    display: flex;
    font-size: 0.85rem;
    color: var(--edna-yellow);
    padding-bottom: 8px;
    border-bottom: 1px solid var(--edna-gray);
    font-weight: bold;
    text-transform: uppercase;
    justify-content: space-between;
}

.items-body {
    flex: 1;
    overflow-y: auto;
    padding-top: 10px;
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
.col-date {
    width: 80px;
    color: var(--edna-light-gray);
    font-size: 0.85rem;
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

/* Total */
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
.debt {
    color: var(--edna-orange);
    text-shadow: 0 0 5px rgba(244, 113, 110, 0.2);
}

/* Botão Fab */
.btn-fab {
    position: absolute;
    bottom: -25px;
    left: 50%;
    transform: translateX(-50%);
    width: 60px;
    height: 60px;
    border-radius: 50%;
    background-color: var(--edna-green);
    color: var(--edna-black);
    font-size: 2.5rem;
    border: 4px solid var(--edna-dark-gray);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.5);
    transition:
        transform 0.2s,
        background-color 0.2s;
    padding-bottom: 6px;
}
.btn-fab:hover {
    transform: translateX(-50%) scale(1.1);
    background-color: var(--edna-white);
}

/* --- PAINEL DIREITO --- */
.panel-right {
    flex: 1;
    padding: 20px;
    background-color: var(--edna-black);
    overflow-y: auto;
}

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
    border-left: 4px solid var(--edna-blue);
}
.card-prod:hover {
    background-color: var(--edna-gray);
    transform: translateY(-5px);
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
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
.card-cliente {
    border-left-color: var(--edna-wine);
}

.new-client-card {
    border-style: dashed;
    border-color: var(--edna-light-gray);
    opacity: 0.8;
    justify-content: center;
}
.new-client-card:hover {
    opacity: 1;
    border-color: var(--edna-green);
    background-color: rgba(90, 211, 176, 0.1);
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
.client-cpf {
    font-size: 0.85rem;
    color: var(--edna-light-gray);
    font-family: monospace;
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

::-webkit-scrollbar {
    width: 8px;
}
::-webkit-scrollbar-track {
    background: var(--edna-black);
}
::-webkit-scrollbar-thumb {
    background: var(--edna-gray);
    border-radius: 4px;
}
::-webkit-scrollbar-thumb:hover {
    background: var(--edna-light-gray);
}
</style>
