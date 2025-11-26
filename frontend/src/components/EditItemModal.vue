<script setup>
import { reactive, watch } from 'vue';

const props = defineProps({
  item: { type: Object, required: true },
  visible: { type: Boolean, default: false }
});

const emit = defineEmits(['close', 'save']);

const form = reactive({
  nome: '',
  categoria: '',
  marca: ''
});

// Preenche o formulário ao abrir
watch(() => props.item, (newVal) => {
  if (newVal && props.visible) {
    form.nome = newVal.nome;
    form.categoria = newVal.categoria;
    form.marca = newVal.marca;
  }
}, { immediate: true });

const fechar = () => emit('close');

const salvar = () => {
  emit('save', {
    id: props.item.id,
    ...form
  });
};
</script>

<template>
  <div v-if="visible" class="modal-overlay" @click.self="fechar">
    <div class="modal-card">
      <div class="modal-header">
        <h3>Editar Item Estrutural</h3>
        <button class="btn-close" @click="fechar">×</button>
      </div>

      <div class="modal-body">
        <div class="form-group">
          <label>Nome do Item</label>
          <input v-model="form.nome" placeholder="Ex: Copo Descartável" />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>Categoria</label>
            <input v-model="form.categoria" placeholder="Ex: Descartáveis" />
          </div>
          <div class="form-group">
            <label>Marca</label>
            <input v-model="form.marca" placeholder="Ex: PlastCop" />
          </div>
        </div>
      </div>

      <div class="modal-footer">
        <button class="btn-cancel" @click="fechar">Cancelar</button>
        <button class="btn-save" @click="salvar">Salvar</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Reutilizando o estilo padrão dos modais E.D.N.A */
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
  border: 1px solid var(--edna-wine); /* Cor Vinho para diferenciar itens estruturais */
  width: 90%;
  max-width: 500px;
  border-radius: 12px;
  padding: 25px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.6);
  color: var(--edna-white);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  border-bottom: 1px solid var(--edna-gray);
  padding-bottom: 10px;
}

.modal-header h3 {
  margin: 0;
  color: var(--edna-wine);
}

.modal-body {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.form-row {
  display: flex;
  gap: 15px;
}

.form-group {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.form-group label {
  font-size: 0.85rem;
  color: var(--edna-light-gray);
}

input {
  background-color: var(--edna-black);
  border: 1px solid var(--edna-gray);
  color: white;
  padding: 10px;
  border-radius: 6px;
  outline: none;
  box-sizing: border-box;
  width: 100%;
}

input:focus {
  border-color: var(--edna-wine);
}

.modal-footer {
  margin-top: 25px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.btn-close {
  background: none;
  border: none;
  color: var(--edna-red);
  font-size: 1.5rem;
  cursor: pointer;
}

.btn-save {
  background-color: var(--edna-green);
  color: var(--edna-black);
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  font-weight: bold;
  cursor: pointer;
}

.btn-cancel {
  background-color: transparent;
  border: 1px solid var(--edna-gray);
  color: var(--edna-light-gray);
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
}

.btn-save:hover { filter: brightness(1.1); }
.btn-cancel:hover { background-color: rgba(255,255,255,0.05); color: white; }

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
