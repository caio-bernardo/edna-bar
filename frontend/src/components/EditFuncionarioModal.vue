<script setup>
import { reactive, watch } from 'vue';

const props = defineProps({
  funcionario: { type: Object, required: true },
  visible: { type: Boolean, default: false }
});

const emit = defineEmits(['close', 'save']);

// Estado local do formulário
const form = reactive({
  nome: '',
  CPF: '',
  tipo: 'garcom',
  expediente: 'noite',
  salario: '',
  data_contratacao: ''
});

// Preenche o formulário ao abrir
watch(() => props.funcionario, (newVal) => {
  if (newVal && props.visible) {
    form.nome = newVal.nome;
    form.CPF = newVal.CPF;
    form.tipo = newVal.tipo;
    form.expediente = newVal.expediente;
    form.salario = newVal.salario;
    
    // Ajuste de data para o input type="date" (YYYY-MM-DD)
    if (newVal.data_contratacao) {
      form.data_contratacao = newVal.data_contratacao.split('T')[0];
    }
  }
}, { immediate: true });

const fechar = () => emit('close');

const salvar = () => {
  // Prepara o payload e emite para o pai
  emit('save', {
    id: props.funcionario.id,
    ...form,
    salario: parseFloat(form.salario)
  });
};
</script>

<template>
  <div v-if="visible" class="modal-overlay" @click.self="fechar">
    <div class="modal-card">
      <div class="modal-header">
        <h3>Editar Funcionário</h3>
        <button class="btn-close" @click="fechar">×</button>
      </div>

      <div class="modal-body">
        <div class="form-group">
          <label>Nome Completo</label>
          <input v-model="form.nome" placeholder="Nome" />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>CPF</label>
            <input v-model="form.CPF" placeholder="000.000.000-00" maxlength="14" />
          </div>
          <div class="form-group">
            <label>Salário (R$)</label>
            <input v-model="form.salario" type="number" step="0.01" />
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>Função</label>
            <select v-model="form.tipo">
              <option value="garcom">Garçom</option>
              <option value="seguranca">Segurança</option>
              <option value="caixa">Caixa</option>
              <option value="faxineiro">Faxineiro</option>
              <option value="balconista">Balconista</option>
            </select>
          </div>
          <div class="form-group">
            <label>Expediente</label>
            <select v-model="form.expediente">
              <option value="noite">Noite</option>
              <option value="manha">Manhã</option>
              <option value="tarde">Tarde</option>
            </select>
          </div>
        </div>

        <div class="form-group">
          <label>Data de Contratação</label>
          <input type="date" v-model="form.data_contratacao" />
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
/* Reutilizando o estilo do modal de produtos para consistência */
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
  border: 1px solid var(--edna-blue); /* Borda azul para diferenciar visualmente */
  width: 90%;
  max-width: 550px;
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
  color: var(--edna-blue);
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

input, select {
  background-color: var(--edna-black);
  border: 1px solid var(--edna-gray);
  color: white;
  padding: 10px;
  border-radius: 6px;
  outline: none;
  box-sizing: border-box;
  width: 100%;
}

input:focus, select:focus {
  border-color: var(--edna-blue);
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
