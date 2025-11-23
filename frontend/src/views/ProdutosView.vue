<script setup>
import { ref, onMounted, reactive } from "vue";
import api from "@/services/api";
import EditProdutoModal from "@/components/EditProdutoModal.vue";

// --- ESTADO ---
const produtos = ref([]);
const ofertas = ref([]);
const carregando = ref(false);

// Estado do Modal de Edição
const showEditModal = ref(false);
const produtoParaEditar = ref({});

// Formulário de Produto
const produtoForm = reactive({
    nome: "",
    marca: "",
    categoria: "",
    preco_venda: "",
});

// Formulário de Oferta
const ofertaForm = reactive({
    nome: "",
    data_inicio: "",
    data_fim: "",
    tipo_valor: "desconto", // 'desconto' ou 'fixo'
    valor: "", // Será mapeado para valor_fixo ou percentual_desconto
});

// --- AÇÕES ---
const carregarDados = async () => {
    carregando.value = true;
    try {
        // Buscando no backend os produtos e ofertas
        const [resProdutos, resOfertas] = await Promise.all([
            api.getProdutosComerciais(),
            api.getOfertas(),
        ]);

        // Enquanto não temos a quantidade
        const produtosTemp = resProdutos.data || [];

        // Busca quantidades em paralelo
        const produtosComQuantidade = await Promise.all(
            produtosTemp.map(async (p) => {
                try {
                    const resQtd = await api.getProdutoQtd(p.id);
                    return {
                        ...p,
                        quantidade: resQtd.data.quantidade_disponível,
                    };
                } catch {
                    return { ...p, quantidade: 0 };
                }
            }),
        );

        produtos.value = produtosComQuantidade;
        ofertas.value = resOfertas.data || [];
    } catch (error) {
        console.error("Erro ao carregar dados:", error);
        alert("Erro ao conectar com o servidor.");
    } finally {
        carregando.value = false;
    }
};

const criarProduto = async () => {
    if (!produtoForm.nome || !produtoForm.preco_venda)
        return alert("Preencha nome e preço.");

    try {
        const payload = {
            nome: produtoForm.nome,
            marca: produtoForm.marca,
            categoria: produtoForm.categoria,
            preco_venda: parseFloat(produtoForm.preco_venda), // Backend espera float32
        };

        await api.createProdutoComercial(payload);

        // Limpar form e recarregar lista
        Object.assign(produtoForm, {
            nome: "",
            marca: "",
            categoria: "",
            preco_venda: "",
        });
        await carregarDados();
    } catch (error) {
        alert(
            "Erro ao criar produto: " +
                (error.response?.data?.detail || error.message),
        );
    }
};

const criarOferta = async () => {
    if (!ofertaForm.nome || !ofertaForm.data_inicio || !ofertaForm.data_fim)
        return alert("Preencha os dados obrigatórios da oferta.");

    try {
        // Preparar datas para o formato que o Go espera (RFC3339: YYYY-MM-DDTHH:MM:SSZ)
        // O input type="date" retorna YYYY-MM-DD, adicionamos o tempo.
        const inicioISO = new Date(ofertaForm.data_inicio).toISOString();
        const fimISO = new Date(ofertaForm.data_fim).toISOString();

        const payload = {
            nome: ofertaForm.nome,
            data_inicio: inicioISO,
            data_fim: fimISO,
            // Envia null nos campos que não estão sendo usados
            valor_fixo:
                ofertaForm.tipo_valor === "fixo"
                    ? parseFloat(ofertaForm.valor)
                    : null,
            percentual_desconto:
                ofertaForm.tipo_valor === "desconto"
                    ? parseInt(ofertaForm.valor)
                    : null,
        };

        await api.createOferta(payload);

        // Limpar form e recarregar
        Object.assign(ofertaForm, {
            nome: "",
            data_inicio: "",
            data_fim: "",
            valor: "",
        });
        await carregarDados();
    } catch (error) {
        alert(
            "Erro ao criar oferta: " +
                (error.response?.data?.detail || error.message),
        );
    }
};

const deletarItem = async (tipo, id) => {
    if (!confirm("Tem certeza que deseja excluir?")) return;
    try {
        const endpoint =
            tipo === "produto" ? `/produtos/${id}` : `/ofertas/${id}`;
        await api.deleteByEndpoint(endpoint);
        await carregarDados();
    } catch (error) {
        alert("Erro ao deletar.");
    }
};

// Abre o modal com os dados do produto clicado
const abrirEdicao = (prod) => {
    produtoParaEditar.value = { ...prod }; // Copia para não alterar a lista diretamente
    showEditModal.value = true;
};

// Recebe os dados salvos do modal e envia para a API
const salvarEdicao = async (dadosAtualizados) => {
    try {
        await api.updateProdutoComercial(dadosAtualizados.id, dadosAtualizados);
        showEditModal.value = false; // Fecha modal
        await carregarDados(); // Atualiza lista
    } catch (error) {
        alert("Erro ao atualizar produto: " + (error.response?.data?.detail || error.message));
    }
};

// Inicialização
onMounted(() => {
    carregarDados();
});
</script>

<template>
    <div class="nav-space"></div>
    <div class="comercial-container">
        <h1 class="page-title">Comercial</h1>

        <EditProdutoModal 
            :visible="showEditModal" 
            :produto="produtoParaEditar" 
            @close="showEditModal = false" 
            @save="salvarEdicao" 
        />

        <div class="content-grid">
            <section class="panel produtos-panel">
                <header>
                    <h2>Produtos</h2>
                    <div class="form-inline">
                        <input
                            v-model="produtoForm.nome"
                            placeholder="Nome"
                            type="text"
                        />
                        <input
                            v-model="produtoForm.marca"
                            placeholder="Marca"
                            type="text"
                        />
                        <input
                            v-model="produtoForm.categoria"
                            placeholder="Categoria"
                            type="text"
                        />
                        <input
                            v-model="produtoForm.preco_venda"
                            placeholder="Preço (R$)"
                            type="number"
                            step="0.01"
                        />
                        <button @click="criarProduto" class="btn-add">
                            <span class="icon">+</span>
                        </button>
                    </div>
                </header>
                <div class="list-cards">
                    <div
                        v-if="produtos.length === 0 && !carregando"
                        class="empty-state"
                    >
                        Nenhum produto.
                    </div>
                    <div v-if="carregando" class="empty-state">
                        Carregando...
                    </div>

                    <div
                        v-for="prod in produtos"
                        :key="prod.id"
                        class="card-item card-produto"
                    >
                        <div class="card-main">
                            <div class="card-top">
                                <h3>{{ prod.nome }}</h3>
                                <span class="marca-tag">{{
                                    prod.marca || "Genérico"
                                }}</span>
                            </div>
                            <div class="card-details">
                                <span class="price"
                                    >R$ {{ prod.preco_venda.toFixed(2) }}</span
                                >
                                <span
                                    class="stock"
                                    :class="{
                                        'low-stock': prod.quantidade <= 5,
                                    }"
                                >
                                    Estoque: {{ prod.quantidade }}
                                </span>
                            </div>
                        </div>
                        <div class="card-actions">
                            <button @click="abrirEdicao(prod)" class="btn-icon-edit" title="Editar">
                                ✏️
                            </button>
                            <button @click="deletarItem('produto', prod.id)" class="btn-icon-delete" title="Excluir">
                                ×
                            </button>
                        </div>
                    </div>
                </div>
            </section>

            <section class="panel ofertas-panel">
                <header>
                    <h2>Ofertas & Promoções</h2>

                    <div class="form-stack">
                        <div class="row">
                            <input
                                v-model="ofertaForm.nome"
                                placeholder="Nome da Promoção (ex: Combo Fritas)"
                            />
                        </div>
                        <div class="row dates">
                            <label
                                >Início:
                                <input
                                    v-model="ofertaForm.data_inicio"
                                    type="date"
                            /></label>
                            <label
                                >Fim:
                                <input
                                    v-model="ofertaForm.data_fim"
                                    type="date"
                            /></label>
                        </div>
                        <div class="row values">
                            <select v-model="ofertaForm.tipo_valor">
                                <option value="desconto">% Desconto</option>
                                <option value="fixo">R$ Fixo</option>
                            </select>
                            <input
                                v-model="ofertaForm.valor"
                                type="number"
                                placeholder="Valor"
                            />
                            <button @click="criarOferta" class="btn-save">
                                Criar
                            </button>
                        </div>
                    </div>
                </header>

                <div class="list-cards">
                    <div v-if="ofertas.length === 0" class="empty-state">
                        Nenhuma oferta ativa.
                    </div>

                    <div
                        v-for="oferta in ofertas"
                        :key="oferta.id_oferta"
                        class="card-oferta"
                    >
                        <div class="card-header">
                            <h3>{{ oferta.nome }}</h3>
                            <button
                                @click="deletarItem('oferta', oferta.id_oferta)"
                                class="btn-icon-delete"
                            >
                                ×
                            </button>
                        </div>

                        <div class="card-body">
                            <p class="dates">
                                📅
                                {{
                                    new Date(
                                        oferta.data_inicio,
                                    ).toLocaleDateString()
                                }}
                                até
                                {{
                                    new Date(
                                        oferta.data_fim,
                                    ).toLocaleDateString()
                                }}
                            </p>
                            <div class="tag-price">
                                <span
                                    v-if="oferta.percentual_desconto"
                                    class="discount"
                                    >-{{ oferta.percentual_desconto }}%
                                    OFF</span
                                >
                                <span v-if="oferta.valor_fixo" class="fixed"
                                    >R$ {{ oferta.valor_fixo }}</span
                                >
                            </div>

                            <div class="itens-placeholder">
                                <small>Itens da oferta:</small>
                                <ul>
                                    <li class="disabled-text">
                                        Lista de itens indisponível (API
                                        pendente)
                                    </li>
                                </ul>
                            </div>
                        </div>
                    </div>
                </div>
            </section>
        </div>
    </div>
</template>

<style scoped>
/* PALETA DE CORES E.D.N.A. */
/* Para syntax highlight colorido :3 */
:root {
    --edna-blue: #b6e5f3;
    --edna-green: #88caaf;
    --edna-wine: #a12d4c;
    --edna-red: #e71d51;
    --edna-orange: #f4716e;
    --edna-yellow: #ffd782;
    --edna-light-gray: #888899;
    --edna-gray: #353548;
    --edna-dark-gray: #2a2a32;
    --edna-black: #1a1a1e;
    --edna-white: #f4f4ff;
}
/* --- VARIÁVEIS GERAIS --- */
.nav-space {
    background-image: linear-gradient(
        220deg,
        var(--edna-orange),
        var(--edna-yellow)
    );
}

.comercial-container {
    background-color: var(--edna-black);
    color: var(--edna-white);
    min-height: 100vh;
    padding: 20px;
    font-family: "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
}

.page-title {
    text-align: center;
    margin-bottom: 30px;
    font-weight: 300;
    font-size: 2rem;
    letter-spacing: 2px;
    border-bottom: 2px solid var(--edna-gray);
    padding-bottom: 10px;
}

/* --- GRID LAYOUT --- */
.content-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 40px;
    align-items: start;
}

@media (max-width: 900px) {
    .content-grid {
        grid-template-columns: 1fr;
    }
}

/* --- PAINÉIS --- */
.panel {
    background-color: var(--edna-dark-gray);
    border: 0px solid var(--edna-gray);
    border-radius: 12px;
    padding: 20px;
    box-shadow: 0 4px 1vw rgba(0, 0, 0, 0.5);
}

.panel h2 {
    margin-top: 0;
    border-bottom: 1px solid var(--edna-gray);
    padding-bottom: 10px;
    margin-bottom: 20px;
    color: var(--edna-yellow);
}

/* --- FORMULÁRIOS --- */
input,
select {
    background-color: var(--edna-gray);
    color: var(--edna-white);
    padding: 8px 12px;
    border-radius: 6px;
    outline: none;
}


/* Estilo do form inline (Produtos) */
.form-inline {
    display: flex;
    gap: 1vw;
    margin-bottom: 3vh;
    flex-wrap: wrap;
}
.form-inline input {
  flex: 1;
  width: auto;
  min-width: 80px;
}

/* Estilo do form stack (Ofertas) */
.form-stack {
    display: flex;
    flex-direction: column;
    gap: 1vw;
    margin-bottom: 3vh;
    padding: 15px;
    border-radius: 8px;
}
.form-stack .row {
    display: flex;
    gap: 1vw;
}
.form-stack .dates label {
    font-size: 0.85rem;
    color: var(--edna-light-gray);
    display: flex;
    flex-direction: column;
}

/* --- BOTÕES --- */
button {
    color: var(--edna-black);
    cursor: pointer;
    border: none;
    border-radius: 6px;
    transition: filter 0.2s;
}
button:hover {
    filter: brightness(1.1);
}

.btn-add {
    background-color: var(--edna-green);
    font-size: 1.5rem;
    width: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.btn-add span {
    color: var(--edna-black);
}

.btn-save {
    background-color: var(--edna-green);
    padding: 0 20px;
}

/* --- CARDS --- */
.card-item {
    background-color: var(--edna-gray);
    border-radius: 8px;
    border: 2px solid rgba(0, 0, 0, 0);
    padding: 15px;
    transition:
        transform 0.1s,
        background-color 0.2s;
}

.card-item:hover {
  border-color: var(--edna-light-gray);
}

/* --- CARDS (PRODUTOS) --- */
.card-produto {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.card-main {
    flex: 1;
}

.card-top {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 6px;
}

.card-top h3 {
    margin: 0;
    font-size: 1.1rem;
    color: var(--edna-white);
}

.marca-tag {
  font-size: 0.8rem;
  color: var(--edna-light-gray);
  text-transform: uppercase;
  background: var(--edna-black);
  padding: 2px 6px;
  border-radius: 4px;
}

.card-details {
    display: flex;
    gap: 15px;
    font-size: 0.9rem;
}

.price {
    color: var(--edna-green);
}

.stock {
    color: var(--edna-light-gray);
}
.low-stock {
    color: var(--edna-orange);
}

.btn-icon-delete {
    background: none;
    color: var(--edna-red);
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    opacity: 0.7;
    padding: 5px;
    transition: opacity 0.2s;
}
.btn-icon-delete:hover {
    transform: scale(1.2);
    opacity: 1;
}

.card-actions {
    display: flex;
    align-items: center;
    gap: 5px;
}

.btn-icon-edit {
    background: none;
    border: none;
    font-size: 1rem;
    cursor: pointer;
    opacity: 0.7;
    transition: transform 0.2s, opacity 0.2s;
    filter: grayscale(100%); 
}

.btn-icon-edit:hover {
    opacity: 1;
    transform: scale(1.2);
    filter: none;
}

/* --- CARDS (OFERTAS) --- */
.list-cards {
    display: flex;
    flex-direction: column;
    gap: 15px;
    max-height: 500px;
    overflow-y: auto;
}

.card-oferta {
    background-color: var(--edna-gray);
    border: 1px solid var(--edna-gray);
    border-radius: 8px;
    padding: 15px;
    position: relative;
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
}

.card-header h3 {
    margin: 0 0 10px 0;
    font-size: 1.1rem;
    color: var(--edna-white);
}

.dates {
    font-size: 0.85rem;
    color: var(--edna-light-gray);
    margin-bottom: 10px;
}

.tag-price span {
    display: inline-block;
    padding: 4px 8px;
    border-radius: 4px;
    font-weight: bold;
    font-size: 0.9rem;
}

.tag-price .discount {
    background-color: var(--edna-orange);
    color: var(--edna-gray);
}

.tag-price .fixed {
    background-color: var(--edna-green);
    color: var(--edna-gray);
}

.itens-placeholder {
    margin-top: 15px;
    padding-top: 10px;
    border-top: 1px dashed var(--edna-gray);
}
.itens-placeholder ul {
    list-style: none;
    padding: 0;
}
.disabled-text {
    color: var(--edna-light-gray);
    font-style: italic;
}

/* SCROLL BAR */
.list-cards::-webkit-scrollbar {
    width: 6px;
}
.list-cards::-webkit-scrollbar-thumb {
    background: #444;
    border-radius: 3px;
}
.list-cards::-webkit-scrollbar-track {
    background: #161616;
}
</style>
