<script setup>
import { reactive, watch, computed } from "vue";

const props = defineProps({
    oferta: { type: Object, required: true },
    produtosDisponiveis: { type: Array, default: () => [] }, // Lista completa de produtos para adicionar
    visible: { type: Boolean, default: false },
});

const emit = defineEmits(["close", "save-info", "add-item", "remove-item"]);

// --- DADOS BÁSICOS DA OFERTA ---
const form = reactive({
    nome: "",
    data_inicio: "",
    data_fim: "",
    tipo_valor: "desconto",
    valor: "",
});

// --- DADOS PARA ADICIONAR NOVO ITEM ---
const formItem = reactive({
    id_produto: "",
    quantidade: 1,
});

// Atualiza o formulário quando a oferta muda
watch(
    () => props.oferta,
    (newVal) => {
        if (newVal && props.visible) {
            form.nome = newVal.nome;

            // Formatar datas para input date (YYYY-MM-DD)
            if (newVal.data_inicio)
                form.data_inicio = newVal.data_inicio.split("T")[0];
            if (newVal.data_fim) form.data_fim = newVal.data_fim.split("T")[0];

            // Determinar tipo de valor (Desconto ou Fixo)
            if (newVal.valor_fixo !== null && newVal.valor_fixo !== undefined) {
                form.tipo_valor = "fixo";
                form.valor = newVal.valor_fixo;
            } else {
                form.tipo_valor = "desconto";
                form.valor = newVal.percentual_desconto;
            }
        }
    },
    { immediate: true },
);

// --- AÇÕES ---

const fechar = () => emit("close");

// 1. Salvar Dados Básicos
const salvarInfo = () => {
    const payload = {
        id: props.oferta.id_oferta,
        nome: form.nome,
        data_inicio: new Date(form.data_inicio).toISOString(),
        data_fim: new Date(form.data_fim).toISOString(),
        valor_fixo: form.tipo_valor === "fixo" ? parseFloat(form.valor) : null,
        percentual_desconto:
            form.tipo_valor === "desconto" ? parseInt(form.valor) : null,
    };
    emit("save-info", payload);
};

// 2. Adicionar Item
const adicionarItem = () => {
    if (!formItem.id_produto || formItem.quantidade < 1)
        return alert("Selecione um produto e quantidade.");

    emit("add-item", {
        id_oferta: props.oferta.id_oferta,
        id_produto: parseInt(formItem.id_produto),
        quantidade: parseInt(formItem.quantidade),
    });

    // Resetar seleção
    formItem.id_produto = "";
    formItem.quantidade = 1;
};

// 3. Remover Item
const removerItem = (item) => {
    if (confirm(`Remover ${item.nomeProduto} desta oferta?`)) {
        emit("remove-item", {
            id_oferta: props.oferta.id_oferta,
            id_produto: item.id_produto,
        });
    }
};
</script>

<template>
    <div v-if="visible" class="modal-overlay" @click.self="fechar">
        <div class="modal-card">
            <div class="modal-header">
                <h3>Editar Promoção</h3>
                <button class="btn-close" @click="fechar">×</button>
            </div>

            <div class="modal-body scrollable">
                <fieldset class="section-box">
                    <legend>Dados da Oferta</legend>
                    <div class="form-group">
                        <label>Nome</label>
                        <input
                            v-model="form.nome"
                            placeholder="Ex: Happy Hour"
                        />
                    </div>
                    <div class="form-row">
                        <div class="form-group">
                            <label>Início</label>
                            <input type="date" v-model="form.data_inicio" />
                        </div>
                        <div class="form-group">
                            <label>Fim</label>
                            <input type="date" v-model="form.data_fim" />
                        </div>
                    </div>
                    <div class="form-row">
                        <div class="form-group">
                            <label>Tipo</label>
                            <select v-model="form.tipo_valor">
                                <option value="desconto">% Desconto</option>
                                <option value="fixo">R$ Fixo</option>
                            </select>
                        </div>
                        <div class="form-group">
                            <label>Valor</label>
                            <input type="number" v-model="form.valor" />
                        </div>
                        <div class="form-group btn-container">
                            <button class="btn-save-info" @click="salvarInfo">
                                Salvar Dados
                            </button>
                        </div>
                    </div>
                </fieldset>

                <fieldset class="section-box">
                    <legend>Itens Inclusos</legend>

                    <div class="add-item-row">
                        <select
                            v-model="formItem.id_produto"
                            class="select-prod"
                        >
                            <option value="" disabled>
                                Selecione um produto...
                            </option>
                            <option
                                v-for="p in produtosDisponiveis"
                                :key="p.id"
                                :value="p.id"
                            >
                                {{ p.nome }} ({{ p.marca }})
                            </option>
                        </select>
                        <input
                            type="number"
                            v-model="formItem.quantidade"
                            min="1"
                            class="input-qtd"
                            placeholder="Qtd"
                        />
                        <button class="btn-add" @click="adicionarItem">
                            +
                        </button>
                    </div>

                    <div class="items-list-container">
                        <ul
                            v-if="oferta.itens && oferta.itens.length > 0"
                            class="items-list"
                        >
                            <li
                                v-for="item in oferta.itens"
                                :key="item.id_produto"
                            >
                                <div class="item-desc">
                                    <span class="qtd"
                                        >{{ item.quantidade }}x</span
                                    >
                                    <span class="name">{{
                                        item.nomeProduto
                                    }}</span>
                                </div>
                                <button
                                    class="btn-remove-item"
                                    @click="removerItem(item)"
                                >
                                    ×
                                </button>
                            </li>
                        </ul>
                        <p v-else class="empty-msg">Nenhum item vinculado.</p>
                    </div>
                </fieldset>
            </div>
        </div>
    </div>
</template>

<style scoped>
.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background-color: rgba(0, 0, 0, 0.7);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
    backdrop-filter: blur(2px);
}

.modal-card {
    background-color: var(--edna-dark-gray);
    border: 1px solid var(--edna-orange);
    width: 95%;
    max-width: 600px;
    max-height: 90vh;
    border-radius: 12px;
    padding: 20px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.6);
    color: var(--edna-white);
    display: flex;
    flex-direction: column;
}

.modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;
    border-bottom: 1px solid var(--edna-gray);
    padding-bottom: 10px;
}
.modal-header h3 {
    margin: 0;
    color: var(--edna-orange);
}

.scrollable {
    overflow-y: auto;
    flex: 1;
    padding-right: 5px;
}

/* Fieldsets */
.section-box {
    border: 1px solid var(--edna-gray);
    border-radius: 8px;
    padding: 15px;
    margin-bottom: 20px;
}
.section-box legend {
    color: var(--edna-light-gray);
    padding: 0 5px;
    font-size: 0.9rem;
}

/* Forms */
.form-group {
    display: flex;
    flex-direction: column;
    gap: 5px;
    margin-bottom: 10px;
    flex: 1;
}
.form-row {
    display: flex;
    gap: 10px;
    align-items: flex-end;
}
.btn-container {
    justify-content: flex-end;
}

input,
select {
    background-color: var(--edna-black);
    border: 1px solid var(--edna-gray);
    color: white;
    padding: 8px;
    border-radius: 4px;
    outline: none;
    width: 100%;
    box-sizing: border-box;
}
input:focus,
select:focus {
    border-color: var(--edna-orange);
}

/* Botões */
.btn-close {
    background: none;
    border: none;
    color: var(--edna-red);
    font-size: 1.5rem;
    cursor: pointer;
}
.btn-save-info {
    background-color: var(--edna-blue);
    color: var(--edna-black);
    border: none;
    padding: 10px;
    border-radius: 4px;
    cursor: pointer;
    font-weight: bold;
    height: 35px;
}
.btn-save-info:hover {
    filter: brightness(1.1);
}

/* Add Item Row */
.add-item-row {
    display: flex;
    gap: 10px;
    margin-bottom: 15px;
}
.select-prod {
    flex: 3;
}
.input-qtd {
    flex: 1;
}
.btn-add {
    background-color: var(--edna-green);
    color: var(--edna-black);
    width: 35px;
    border-radius: 4px;
    border: none;
    font-size: 1.2rem;
    cursor: pointer;
}

/* Lista de Itens */
.items-list {
    list-style: none;
    padding: 0;
    margin: 0;
}
.items-list li {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px;
    background-color: rgba(255, 255, 255, 0.05);
    margin-bottom: 5px;
    border-radius: 4px;
}
.item-desc {
    display: flex;
    gap: 10px;
}
.qtd {
    color: var(--edna-yellow);
    font-weight: bold;
}
.name {
    color: var(--edna-white);
}

.btn-remove-item {
    background: none;
    border: none;
    color: var(--edna-red);
    font-size: 1.2rem;
    cursor: pointer;
}
.btn-remove-item:hover {
    transform: scale(1.2);
}

.empty-msg {
    color: var(--edna-light-gray);
    font-style: italic;
    text-align: center;
}

/* Scrollbar */
::-webkit-scrollbar {
    width: 6px;
}
::-webkit-scrollbar-thumb {
    background: var(--edna-gray);
    border-radius: 3px;
}

@media (max-width: 600px) {
  .modal-card {
    width: 95% !important;
    max-width: 95% !important;
    padding: 15px;
    max-height: 90vh;
    overflow-y: auto;
  }

  .form-row {
    flex-direction: column;
    gap: 10px;
  }
}
</style>
