<script setup>
import { reactive, watch } from 'vue';

const props = defineProps({
  lote: { type: Object, required: true },
  visible: { type: Boolean, default: false },
  fornecedores: { type: Array, default: () => [] },
  produtos: { type: Array, default: () => [] }
});

const emit = defineEmits(['close', 'save']);

const form = reactive({
  id_fornecedor: '',
  id_produto: '',
  quantidade_inicial: 0,
  preco_unitario: 0,
  estragados: 0,
  data_fornecimento: '',
  validade: ''
});

// Preenche o formulário ao abrir
watch(() => props.lote, (newVal) => {
  if (newVal && props.visible) {
    form.id_fornecedor = newVal.id_fornecedor;
    form.id_produto = newVal.id_produto;
    form.quantidade_inicial = newVal.quantidade_inicial;
    form.preco_unitario = newVal.preco_unitario;
    form.estragados = newVal.estragados || 0;
    
    // Ajuste de datas para inputs type="date"
    if (newVal.data_fornecimento) form.data_fornecimento = newVal.data_fornecimento.split('T')[0];
    if (newVal.validade) form.validade = newVal.validade.split('T')[0];
    else form.validade = '';
  }
}, { immediate: true });

const fechar = () => emit('close');

const salvar = () => {
  // Prepara payload convertendo tipos
  const payload = {
    id: props.lote.id, // ou id_lote dependendo de como vem do prop
    id_fornecedor: form.id_fornecedor,
    id_produto: form.id_produto,
    quantidade_inicial: parseInt(form.quantidade_inicial),
    preco_unitario: parseFloat(form.preco_unitario),
    estragados: parseInt(form.estragados),
    data_fornecimento: new Date(form.data_fornecimento).toISOString(),
    validade: form.validade ? new Date(form.validade).toISOString() : null
  };
  emit('save', payload);
};
</script>

<template>
  <div v-if="visible" class="modal-overlay" @click.self="fechar">
    <div class="modal-card">
      <div class="modal-header">
        <h3>Editar Lote</h3>
        <button class="btn-close" @click="fechar">×</button>
      </div>

      <div class="modal-body">
        
        <div class="form-row">
          <div class="form-group">
            <label>Produto</label>
            <select v-model="form.id_produto">
              <option v-for="p in produtos" :key="p.id" :value="p.id">{{ p.nome }}</option>
            </select>
          </div>
          <div class="form-group">
            <label>Fornecedor</label>
            <select v-model="form.id_fornecedor">
              <option v-for="f in fornecedores" :key="f.id" :value="f.id">{{ f.nome }}</option>
            </select>
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>Qtd Inicial</label>
            <input v-model="form.quantidade_inicial" type="number" />
          </div>
          <div class="form-group">
            <label>Preço Unit. (R$)</label>
            <input v-model="form.preco_unitario" type="number" step="0.01" />
          </div>
          <div class="form-group">
            <label>Itens Estragados</label>
            <input v-model="form.estragados" type="number" />
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>Data Entrada</label>
            <input v-model="form.data_fornecimento" type="date" />
          </div>
          <div class="form-group">
            <label>Validade</label>
            <input v-model="form.validade" type="date" />
          </div>
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
.modal-overlay {
  position: fixed; top: 0; left: 0; width: 100vw; height: 100vh;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex; justify-content: center; align-items: center;
  z-index: 1000; backdrop-filter: blur(2px);
}

.modal-card {
  background-color: var(--edna-dark-gray);
  border: 1px solid var(--edna-blue);
  width: 90%; max-width: 600px;
  border-radius: 12px; padding: 25px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.6);
  color: var(--edna-white);
}

.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  margin-bottom: 20px; border-bottom: 1px solid var(--edna-gray); padding-bottom: 10px;
}
.modal-header h3 { margin: 0; color: var(--edna-blue); }

.modal-body { display: flex; flex-direction: column; gap: 15px; }
.form-row { display: flex; gap: 15px; }
.form-group { flex: 1; display: flex; flex-direction: column; gap: 5px; }
.form-group label { font-size: 0.85rem; color: var(--edna-light-gray); }

input, select {
  background-color: var(--edna-black); border: 1px solid var(--edna-gray);
  color: white; padding: 10px; border-radius: 6px; outline: none;
  box-sizing: border-box; width: 100%;
}
input:focus, select:focus { border-color: var(--edna-blue); }

.modal-footer {
  margin-top: 25px; display: flex; justify-content: flex-end; gap: 10px;
}

.btn-close { background: none; border: none; color: var(--edna-red); font-size: 1.5rem; cursor: pointer; }
.btn-save {
  background-color: var(--edna-green); color: var(--edna-black); border: none;
  padding: 10px 20px; border-radius: 6px; font-weight: bold; cursor: pointer;
}
.btn-cancel {
  background-color: transparent; border: 1px solid var(--edna-gray); color: var(--edna-light-gray);
  padding: 10px 20px; border-radius: 6px; cursor: pointer;
}
.btn-save:hover { filter: brightness(1.1); }
.btn-cancel:hover { background-color: rgba(255,255,255,0.05); color: white; }
</style>
