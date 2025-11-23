<script setup>
import { reactive, watch } from "vue";

const props = defineProps({
    oferta: { type: Object, required: true },
    produtosDisponiveis: { type: Array, default: () => [] },
    visible: { type: Boolean, default: false },
});

const emit = defineEmits(["close", "save-info", "add-item", "remove-item"]);

// --- DADOS BÁSICOS ---
const form = reactive({
    nome: "",
    data_inicio: "",
    data_fim: "",
    tipo_valor: "desconto",
    valor: "",
});

// --- NOVO ITEM ---
const formItem = reactive({
    id_produto: "",
    quantidade: 1,
});

watch(
    () => props.oferta,
    (newVal) => {
        if (newVal && props.visible) {
            form.nome = newVal.nome;
            if (newVal.data_inicio)
                form.data_inicio = newVal.data_inicio.split("T")[0];
            if (newVal.data_fim) form.data_fim = newVal.data_fim.split("T")[0];

            if (newVal.valor_fixo) {
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

const fechar = () => emit("close");

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

const adicionarItem = () => {
    if (!formItem.id_produto) return alert("Selecione um produto.");
    emit("add-item", {
        id_oferta: props.oferta.id_oferta,
        id_produto: parseInt(formItem.id_produto),
        quantidade: parseInt(formItem.quantidade),
    });
    formItem.id_produto = "";
    formItem.quantidade = 1;
};

const removerItem = (item) => {
    if (confirm("Remover item?")) {
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

            <div class="modal-body">
                <div class="section-box">
                    <h4>Dados</h4>
                    <div class="form-row">
                        <input
                            v-model="form.nome"
                            placeholder="Nome"
                            class="input-full"
                        />
                    </div>
                    <div class="form-row">
                        <input type="date" v-model="form.data_inicio" />
                        <input type="date" v-model="form.data_fim" />
                    </div>
                    <div class="form-row">
                        <select v-model="form.tipo_valor">
                            <option value="desconto">% OFF</option>
                            <option value="fixo">R$ Fixo</option>
                        </select>
                        <input
                            type="number"
                            v-model="form.valor"
                            placeholder="Valor"
                        />
                        <button class="btn-save-info" @click="salvarInfo">
                            Salvar
                        </button>
                    </div>
                </div>

                <div class="section-box">
                    <h4>Itens</h4>
                    <div class="add-row">
                        <select v-model="formItem.id_produto">
                            <option value="" disabled>Produto...</option>
                            <option
                                v-for="p in produtosDisponiveis"
                                :key="p.id"
                                :value="p.id"
                            >
                                {{ p.nome }}
                            </option>
                        </select>
                        <input
                            type="number"
                            v-model="formItem.quantidade"
                            min="1"
                            class="input-qtd"
                        />
                        <button class="btn-add" @click="adicionarItem">
                            +
                        </button>
                    </div>
                    <ul class="items-list">
                        <li v-for="item in oferta.itens" :key="item.id_produto">
                            <span
                                >{{ item.quantidade }}x
                                {{ item.nomeProduto }}</span
                            >
                            <button
                                class="btn-remove"
                                @click="removerItem(item)"
                            >
                                ×
                            </button>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.8);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 999;
}
.modal-card {
    background: var(--edna-dark-gray);
    border: 1px solid var(--edna-orange);
    padding: 20px;
    border-radius: 10px;
    width: 90%;
    max-width: 500px;
    color: white;
}
.modal-header {
    display: flex;
    justify-content: space-between;
    margin-bottom: 15px;
    border-bottom: 1px solid #444;
    padding-bottom: 10px;
}
.btn-close {
    background: none;
    border: none;
    color: var(--edna-red);
    font-size: 1.5rem;
    cursor: pointer;
}
.section-box {
    margin-bottom: 20px;
    border: 1px solid #444;
    padding: 10px;
    border-radius: 6px;
}
.section-box h4 {
    margin-top: 0;
    color: var(--edna-light-gray);
    font-size: 0.9rem;
}
.form-row,
.add-row {
    display: flex;
    gap: 10px;
    margin-bottom: 10px;
}
.input-full {
    width: 100%;
}
.input-qtd {
    width: 60px;
}
.btn-save-info {
    background: var(--edna-blue);
    color: black;
    border: none;
    padding: 0 15px;
    border-radius: 4px;
    cursor: pointer;
    font-weight: bold;
}
.btn-add {
    background: var(--edna-green);
    color: black;
    border: none;
    width: 30px;
    border-radius: 4px;
    cursor: pointer;
}
.items-list {
    list-style: none;
    padding: 0;
}
.items-list li {
    display: flex;
    justify-content: space-between;
    padding: 5px;
    background: rgba(255, 255, 255, 0.05);
    margin-bottom: 5px;
    border-radius: 4px;
}
.btn-remove {
    background: none;
    border: none;
    color: var(--edna-red);
    cursor: pointer;
}
</style>
