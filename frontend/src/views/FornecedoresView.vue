<script setup>
import { ref, onMounted, reactive, computed } from 'vue';
import api from '@/services/api';
import EditLoteModal from '@/components/EditLoteModal.vue';

// --- ESTADO ---
const fornecedores = ref([]);
const lotes = ref([]);
const produtos = ref([]);
const carregando = ref(false);

// Estados do Modal de Edição de lote
const showEditLoteModal = ref(false);
const loteParaEditar = ref({});

// Formulário de Fornecedor
const formFornecedor = reactive({
  nome: '',
  cnpj: ''
});

// Formulário de Lote
const formLote = reactive({
  id_fornecedor: '',
  id_produto: '',
  quantidade_inicial: '',
  preco_unitario: '',
  data_fornecimento: '',
  validade: '',
  estragados: 0
});

// --- HELPERS DE EXIBIÇÃO ---
// Como o Lote vem apenas com IDs, precisamos encontrar os nomes nas listas carregadas
const getNomeProduto = (id) => {
  const p = produtos.value.find(item => item.id === id);
  return p ? p.nome : `Produto #${id}`;
};

const getNomeFornecedor = (id) => {
  const f = fornecedores.value.find(item => item.id === id);
  return f ? f.nome : `Fornecedor #${id}`;
};

const formatarData = (dataISO) => {
  if (!dataISO) return 'N/A';
  return new Date(dataISO).toLocaleDateString('pt-BR');
};

// --- AÇÕES ---
const carregarDados = async () => {
  carregando.value = true;
  try {
    const [resFornecedores, resLotes, resProdutos] = await Promise.all([
      api.getFornecedores(),
      api.getLotes(),
      api.getProdutos()
    ]);

    fornecedores.value = resFornecedores.data || [];
    lotes.value = resLotes.data || [];
    produtos.value = resProdutos.data || [];
  } catch (error) {
    console.error("Erro ao carregar:", error);
    alert("Erro de conexão com a API.");
  } finally {
    carregando.value = false;
  }
};

// Criar Fornecedor
const criarFornecedor = async () => {
  if (!formFornecedor.nome || !formFornecedor.cnpj) return alert("Preencha todos os campos.");
  try {
    await api.createFornecedor({ ...formFornecedor });
    Object.assign(formFornecedor, { nome: '', cnpj: '' });
    await carregarDados();
  } catch (error) {
    alert("Erro ao criar fornecedor.");
  }
};

// Criar Lote
const criarLote = async () => {
  // Validação simples
  if (!formLote.id_fornecedor || !formLote.id_produto || !formLote.data_fornecimento) {
    return alert("Preencha os dados obrigatórios do lote.");
  }

  try {
    // Formatar datas para RFC3339 se necessário, ou YYYY-MM-DD dependendo do backend
    // O backend Go espera ISO normalmente em JSON
    const payload = {
      id_fornecedor: formLote.id_fornecedor,
      id_produto: formLote.id_produto,
      quantidade_inicial: parseInt(formLote.quantidade_inicial),
      preco_unitario: parseFloat(formLote.preco_unitario),
      estragados: parseInt(formLote.estragados || 0),
      data_fornecimento: new Date(formLote.data_fornecimento).toISOString(),
      validade: formLote.validade ? new Date(formLote.validade).toISOString() : null
    };

    await api.createLote(payload);
    
    // Resetar form
    Object.assign(formLote, { 
      id_fornecedor: '', id_produto: '', quantidade_inicial: '', 
      preco_unitario: '', data_fornecimento: '', validade: '', estragados: 0 
    });
    
    await carregarDados();
  } catch (error) {
    alert("Erro ao registrar lote: " + (error.response?.data?.detail || error.message));
  }
};

// Deleção
const deletarItem = async (tipo, id) => {
  if(!confirm("Tem certeza? Isso pode afetar o histórico.")) return;
  try {
    if (tipo === 'fornecedor') await api.deleteFornecedor(id);
    if (tipo === 'lote') await api.deleteLote(id);
    await carregarDados();
  } catch (error) {
    alert("Erro ao deletar item.");
  }
};

const abrirEdicaoLote = (lote) => {
  loteParaEditar.value = { ...lote }; // Clona para não alterar a lista
  showEditLoteModal.value = true;
};

const salvarEdicaoLote = async (dados) => {
  try {
    // O ID vem como 'id_lote' do objeto lote original no banco
    await api.updateLote(loteParaEditar.value.id_lote, dados);
    showEditLoteModal.value = false;
    await carregarDados();
  } catch (error) {
    alert("Erro ao atualizar lote: " + (error.response?.data?.detail || error.message));
  }
};

onMounted(() => {
  carregarDados();
});
</script>

<template>
  <div class="nav-space"></div>
  <div class="page-container">
    
    <div class="palette-provider">
      <h1 class="page-title">Estoque & Fornecedores</h1>

      <EditLoteModal 
      :visible="showEditLoteModal"
      :lote="loteParaEditar"
      :fornecedores="fornecedores"
      :produtos="produtos"
      @close="showEditLoteModal = false"
      @save="salvarEdicaoLote"
      />

      <div class="grid-layout">
        
        <section class="panel">
          <header>
            <h2>Lotes de Entrada</h2>
            
            <div class="form-stack">
              <div class="form-row">
                <select v-model="formLote.id_fornecedor" class="flex-2">
                  <option value="" disabled selected>Fornecedor...</option>
                  <option v-for="f in fornecedores" :key="f.id" :value="f.id">{{ f.nome }}</option>
                </select>
                <select v-model="formLote.id_produto" class="flex-2">
                  <option value="" disabled selected>Produto...</option>
                  <option v-for="p in produtos" :key="p.id" :value="p.id">{{ p.nome }}</option>
                </select>
              </div>

              <div class="form-row">
                <input v-model="formLote.quantidade_inicial" type="number" placeholder="Qtd" class="input-small">
                <input v-model="formLote.preco_unitario" type="number" step="0.01" placeholder="R$ Unit." class="input-small">
                <input v-model="formLote.estragados" type="number" placeholder="Estragados" class="input-small">
              </div>

              <div class="form-row">
                <label class="date-label">Entrada: <input v-model="formLote.data_fornecimento" type="date"></label>
                <label class="date-label">Validade: <input v-model="formLote.validade" type="date"></label>
                <button @click="criarLote" class="btn-add btn-full">Registrar</button>
              </div>
            </div>
          </header>

          <div class="scroll-list">
            <div v-if="lotes.length === 0" class="empty-state">Nenhum lote registrado.</div>
            
            <div v-for="lote in lotes" :key="lote.id_lote" class="card lote-card">
              <div class="card-header">
                <h3 class="product-name">{{ getNomeProduto(lote.id_produto) }}</h3>
                <span class="supplier-tag">{{ getNomeFornecedor(lote.id_fornecedor) }}</span>
              </div>

              <div class="card-body">
                <div class="info-grid">
                  <div class="info-item">
                    <span class="label">Qtd Inicial</span>
                    <span class="value">{{ lote.quantidade_inicial }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">Custo Total</span>
                    <span class="value money">R$ {{ (lote.preco_unitario * lote.quantidade_inicial).toFixed(2) }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">Entrada</span>
                    <span class="value">{{ formatarData(lote.data_fornecimento) }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">Validade</span>
                    <span class="value" :class="{'expired': false}">{{ formatarData(lote.validade) }}</span>
                  </div>
                </div>
                
                <div class="card-footer">
                  <span v-if="lote.estragados > 0" class="bad-items">{{ lote.estragados }} estragados</span>
                  <span v-else class="ok-items">0 estragados</span>
                  <div class="actions-group">
                    <button @click="abrirEdicaoLote(lote)" class="btn-icon-edit" title="Editar">✏️</button>
                    <button @click="deletarItem('lote', lote.id_lote)" class="btn-icon-delete" title="Excluir">×</button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section class="panel">
          <header>
            <h2>Fornecedores</h2>
            <div class="form-inline">
              <input v-model="formFornecedor.nome" placeholder="Nome da Empresa" class="flex-2">
              <input v-model="formFornecedor.cnpj" placeholder="CNPJ" class="flex-1">
              <button @click="criarFornecedor" class="btn-icon-add">+</button>
            </div>
          </header>

          <div class="scroll-list">
            <div v-if="fornecedores.length === 0" class="empty-state">Nenhum fornecedor.</div>

            <div v-for="f in fornecedores" :key="f.id" class="card fornecedor-card">
              <div class="f-info">
                <h3>{{ f.nome }}</h3>
                <p class="cnpj">{{ f.cnpj }}</p>
              </div>
              <div class="f-actions">
                <button @click="deletarItem('fornecedor', f.id)" class="btn-icon-delete">×</button>
              </div>
            </div>
          </div>
        </section>

      </div>
    </div>
  </div>
</template>

<style scoped>
/* --- PALETA DE CORES E VARIÁVEIS --- */
.page-container {
  /* Definindo as cores da imagem fornecida */
  --edna-blue: #B6E5F3;
  --edna-green: #5AD3B0;
  --edna-wine: #A12D4C;
  --edna-red: #E71D51;
  --edna-orange: #F4716E;
  --edna-yellow: #FFD782;
  --edna-light-gray: #888899;
  --edna-gray: #353545;
  --edna-dark-gray: #2A2A32;
  --edna-black: #1A1A1E;
  --edna-white: #F4F4FF;

  background-color: var(--edna-black);
  color: var(--edna-white);
  height: 100vh; 
  display: flex;
  flex-direction: column;
  overflow: hidden; /* Impede scroll na página inteira */
  padding: 20px;
  box-sizing: border-box;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

/* --- LAYOUT GERAL --- */
.nav-space {
  background-image: linear-gradient(220deg, var(--edna-red), var(--edna-orange));
}

.grid-layout {
  display: grid;
  grid-template-columns: 3fr 2fr; /* Esquerda maior (Lotes), Direita menor */
  gap: 30px;
  height: calc(100vh - 120px); /* Tenta ocupar a tela toda */
}

@media (max-width: 900px) {
  .grid-layout {
    grid-template-columns: 1fr;
    height: auto;
  }
}

/* --- PAINÉIS --- */
.panel {
  background-color: var(--edna-dark-gray);
  border: 1px solid var(--edna-gray);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  box-shadow: 0 4px 6px rgba(0,0,0,0.3);
  box-sizing: border-box;
}

.panel header {
  margin-bottom: 15px;
  border-bottom: 1px solid var(--edna-gray);
  padding-bottom: 15px;
  flex-shrink: 0;
}

/* --- SCROLL AREA --- */
.scroll-list {
  flex: 1;
  overflow-y: auto;
  padding-right: 8px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-height: 0;
}

.scroll-list::-webkit-scrollbar { width: 6px; }
.scroll-list::-webkit-scrollbar-thumb { background: var(--edna-gray); border-radius: 3px; }
.scroll-list::-webkit-scrollbar-track { background: var(--edna-black); }

/* --- FORMS --- */
.form-stack { display: flex; flex-direction: column; gap: 10px; }
.form-row, .form-inline { display: flex; gap: 10px; width: 100%; }
.flex-1 { flex: 1; }
.flex-2 { flex: 2; }
.input-small { width: 100px; }

.input, select {
  box-sizing: border-box; /* Mantém o padding dentro da largura */
  min-width: 0; /* Permite que o input encolha se o espaço for pequeno */
  width: 100%;
}
.form-inline { 
  display: flex; 
  gap: 10px; 
  width: 100%; /* Garante que a linha ocupe a largura do painel */
  align-items: center; /* Centraliza verticalmente com o botão */
}
.flex-1 { 
  flex: 1; 
  min-width: 0; /* Reforço para o item flexível */
}

.flex-2 { 
  flex: 2; 
  min-width: 0; /* Reforço para o item flexível */
}

.date-label {
  display: flex;
  flex-direction: column;
  font-size: 0.75rem;
  color: var(--edna-light-gray);
  flex: 1;
}

/* Botões */
.btn-add {
  background-color: var(--edna-green);
  color: var(--edna-black);
  padding: 0 15px;
}
.btn-full { width: 40%; margin-top: auto; height: 2rem; } /* Alinha com os inputs */

.btn-icon-add {
  background-color: var(--edna-green);
  color: var(--edna-black);
  width: 2rem;
  aspect-ratio: 1;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0; 
  flex-grow: 0;
}

.btn-icon {
  background: transparent;
  font-size: 1.1rem;
  padding: 5px;
  opacity: 0.7;
}

.btn-icon:hover {
  opacity: 1;
  transform: scale(1.1);
}

.actions-group {
  display: flex;
  gap: 10px;
  align-items: center;
}

.btn-icon-edit {
  background: none;
  border: none;
  font-size: 1.1rem;
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

/* --- CARDS --- */
.card {
  background-color: var(--edna-gray);
  border-left: 5px solid var(--edna-light-gray);
  border-radius: 5px;
  padding: 12px;
}

.card:hover {
  border-left-color: var(--edna-white);
}

/* Card Lote */
.lote-card {
  border-left-color: var(--edna-orange);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.product-name {
  margin: 0;
  font-size: 1.1rem;
  color: white;
}

.supplier-tag { 
  font-size: 0.9rem; 
  background: var(--edna-black); 
  color: var(--edna-light-gray);
  padding: 2px 6px; 
  border-radius: 4px; 
}

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr 1fr;
  gap: 10px;
  margin-bottom: 10px;
}

.info-item {
  display: flex;
  flex-direction: column;
}

.info-item .label {
  font-size: 0.9rem;
  color: var(--edna-light-gray); 
  text-transform: uppercase;
  margin-bottom: 0.2rem;
}

.info-item .value {
  font-size: 1rem;
}

.money { color: var(--edna-green); }

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px dashed var(--edna-gray);
  font-size: 0.85rem;
}

.bad-items { color: var(--edna-orange);}
.ok-items { color: var(--edna-green);}

/* Card Fornecedor */
.fornecedor-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-left-color: var(--edna-red);
}
.f-info h3 {
  margin: 0;
  font-size: 1.1rem;
  color: white;
}
.f-info .cnpj {
  margin: 2px 0 0 0;
  font-size: 0.9rem;
  color: var(--edna-light-gray); }

.empty-state {
  text-align: center;
  color: var(--edna-light-gray);
  margin-top: 20px;
  font-style: italic;
}

/* --- MOBILE --- */
@media (max-width: 900px) {
  .page-container {
    height: auto;
    min-height: 100vh;
    overflow-y: auto;
    display: block;
  }

  .grid-layout {
    display: flex;
    flex-direction: column;
    height: auto;
    overflow: visible;
  }
  
  .panel {
    /* Altura fixa ou mínima para cada painel no mobile */
    height: 600px;
    min-height: 0;
    margin-bottom: 30px;
    overflow: hidden;
  }
  
  .page-container::-webkit-scrollbar {
    width: 8px;
  }
}
</style>
