<script setup>
import { ref } from "vue";

const props = defineProps({
    venda: {
        type: Object,
        required: true,
    },
});

const emit = defineEmits(["close", "confirm"]);

const tipoPagamento = ref("pix");
const processando = ref(false);

const confirmar = () => {
    processando.value = true;
    emit("confirm", {
        ...props.venda,
        tipo_pagamento: tipoPagamento.value,
    });
};
</script>

<template>
    <div class="modal-backdrop" @click.self="$emit('close')">
        <div class="modal-card">
            <div class="modal-header">
                <h3>Pagar Conta (Fiado)</h3>
                <button class="btn-close" @click="$emit('close')">×</button>
            </div>

            <div class="modal-body">
                <div class="info-row">
                    <span>Venda ID:</span>
                    <strong>#{{ venda.id }}</strong>
                </div>
                <div class="info-row">
                    <span>Valor Total:</span>
                    <span class="valor">R$ {{ venda.total.toFixed(2) }}</span>
                </div>
                <div class="info-row">
                    <span>Data da Compra:</span>
                    <span>{{
                        new Date(venda.data_hora_renda).toLocaleDateString(
                            "pt-BR",
                        )
                    }}</span>
                </div>

                <hr class="divider" />

                <div class="form-group">
                    <label>Forma de Pagamento:</label>
                    <select v-model="tipoPagamento">
                        <option value="pix">Pix</option>
                        <option value="credito">Crédito</option>
                        <option value="debito">Débito</option>
                        <option value="dinheiro">Dinheiro</option>
                    </select>
                </div>
            </div>

            <div class="modal-footer">
                <button class="btn-cancel" @click="$emit('close')">
                    Cancelar
                </button>
                <button
                    class="btn-confirm"
                    @click="confirmar"
                    :disabled="processando"
                >
                    {{ processando ? "Processando..." : "Confirmar Pagamento" }}
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.modal-backdrop {
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
}

.modal-card {
    background-color: var(--edna-dark-gray);
    border: 1px solid var(--edna-gray);
    border-radius: 12px;
    width: 90%;
    max-width: 400px;
    color: var(--edna-white);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
    animation: fadeIn 0.2s ease-out;
}

.modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 15px 20px;
    border-bottom: 1px solid var(--edna-gray);
}

.modal-header h3 {
    margin: 0;
    color: var(--edna-yellow);
}

.btn-close {
    background: none;
    border: none;
    color: var(--edna-light-gray);
    font-size: 1.5rem;
    cursor: pointer;
}

.modal-body {
    padding: 20px;
}

.info-row {
    display: flex;
    justify-content: space-between;
    margin-bottom: 10px;
    font-size: 0.9rem;
}

.valor {
    color: var(--edna-green);
    font-weight: bold;
    font-size: 1.1rem;
}

.divider {
    border: 0;
    border-top: 1px dashed var(--edna-gray);
    margin: 15px 0;
}

.form-group {
    display: flex;
    flex-direction: column;
    gap: 5px;
}

select {
    background-color: var(--edna-black);
    color: var(--edna-white);
    border: 1px solid var(--edna-gray);
    padding: 10px;
    border-radius: 6px;
}

.modal-footer {
    padding: 15px 20px;
    background-color: rgba(0, 0, 0, 0.2);
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    border-bottom-left-radius: 12px;
    border-bottom-right-radius: 12px;
}

.btn-cancel {
    background: transparent;
    color: var(--edna-light-gray);
    padding: 8px 16px;
}

.btn-confirm {
    background-color: var(--edna-green);
    color: var(--edna-dark-gray);
    padding: 8px 16px;
    font-weight: bold;
}

.btn-confirm:hover {
    filter: brightness(1.1);
}

@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(-10px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
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
