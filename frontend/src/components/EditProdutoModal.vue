<script setup>
import { reactive, watch } from 'vue';

const props = defineProps({
  produto: { type: Object, required: true }, // Produto original para editar
  visible: { type: Boolean, default: false } // Controle de visibilidade
});

const emit = defineEmits(['close', 'save']);

// Estado local do formulário
const form = reactive({
  nome: '',
  marca: '',
  categoria: '',
  preco_venda: ''
});

// Preenche o formulário sempre que o produto muda ou o modal abre
watch(() => props.produto, (newVal) => {
  if (newVal) {
    form.nome = newVal.nome;
    form.marca = newVal.marca;
    form.categoria = newVal.categoria;
    form.preco_venda = newVal.preco_venda;
  }
}, { immediate: true });

const fechar = () => {
  emit('close');
};

const salvar = () => {
  // Emite os dados editados de volta para a View pai
  emit('save', {
    id: props.produto.id,
    ...form,
    preco_venda: parseFloat(form.preco_venda)
  });
};
</script>

<template>
  <div v-if="visible" class="modal-overlay" @click.self="fechar">
    <div class="modal-card">
      <div class="modal-header">
        <h3>Editar Produto</h3>
        <button class="btn-close" @click="fechar">×</button>
      </div>

      <div class="modal-body">
        <div class="form-group">
          <label>Nome</label>
          <input v-model="form.nome" placeholder="Nome do Produto" />
        </div>
        
        <div class="form-row">
          <div class="form-group">
            <label>Marca</label>
            <input v-model="form.marca" placeholder="Marca" />
          </div>
          <div class="form-group">
            <label>Categoria</label>
            <input v-model="form.categoria" placeholder="Categoria" />
          </div>
        </div>

        <div class="form-group">
          <label>Preço de Venda (R$)</label>
          <input v-model="form.preco_venda" type="number" step="0.01" placeholder="0.00" />
        </div>
      </div>

      <div class="modal-footer">
        <button class="btn-cancel" @click="fechar">Cancelar</button>
        <button class="btn-save" @click="salvar">Salvar Alterações</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Fundo Escuro Transparente */
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
  z-index: 1000; /* Fica na frente de tudo */
  backdrop-filter: blur(2px);
}

/* O Cartão do Modal */
.modal-card {
  background-color: var(--edna-dark-gray);
  border: 1px solid var(--edna-yellow);
  width: 90%;
  max-width: 500px;
  border-radius: 12px;
  padding: 25px;
  box-shadow: 0 10px 25px rgba(0,0,0,0.5);
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
  color: var(--edna-yellow);
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

.form-row {
  display: flex;
  gap: 15px;
  width: 100%; 
}

.form-group {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 5px;
  min-width: 0; 
}

input {
  background-color: var(--edna-black);
  border: 1px solid var(--edna-gray);
  color: white;
  padding: 10px;
  border-radius: 6px;
  outline: none;
  width: 100%;
  box-sizing: border-box;
}

input:focus {
  border-color: var(--edna-blue);
}

.modal-footer {
  margin-top: 25px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

/* Botões */
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
</style>
