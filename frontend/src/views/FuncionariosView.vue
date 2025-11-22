<script setup>
import { ref, onMounted, reactive } from 'vue';
import api from '@/services/api';

// --- ESTADO ---
const funcionarios = ref([]);
const itensEstrutura = ref([]);
const carregando = ref(false);

// Formulário Funcionário (Baseado no Model do Swagger/Go)
const formFuncionario = reactive({
  nome: '',
  CPF: '',
  tipo: 'garcom', // Default
  expediente: 'noite', // Default
  salario: '',
  data_contratacao: new Date().toISOString().split('T')[0] // Hoje YYYY-MM-DD
});

// Formulário Item Estrutural
const formItem = reactive({
  nome: '',
  categoria: '',
  marca: ''
});

// --- AÇÕES ---

const carregarDados = async () => {
  carregando.value = true;
  try {
    const [resFunc, resEst] = await Promise.all([
      api.getFuncionarios(),
      api.getProdutosEstruturais()
    ]);

    funcionarios.value = resFunc.data || [];

    // Carregar quantidades dos itens estruturais (opcional, se necessário)
    const itensTemp = resEst.data || [];
    itensEstrutura.value = await Promise.all(itensTemp.map(async (p) => {
      try {
        const resQtd = await api.getProdutoQtd(p.id);
        return { ...p, quantidade: resQtd.data.quantidade_disponível };
      } catch {
        return { ...p, quantidade: 0 };
      }
    }));

  } catch (error) {
    console.error("Erro ao carregar:", error);
  } finally {
    carregando.value = false;
  }
};

const salvarFuncionario = async () => {
  if(!formFuncionario.nome || !formFuncionario.CPF) return alert("Dados incompletos");

  try {
    // Salario deve ser float no Go
    const payload = {
      ...formFuncionario,
      salario: parseFloat(formFuncionario.salario)
    };

    await api.createFuncionario(payload);

    // Resetar form
    Object.assign(formFuncionario, {
      nome: '', CPF: '', salario: '',
      tipo: 'garcom', expediente: 'noite',
      data_contratacao: new Date().toISOString().split('T')[0]
    });

    await carregarDados();
  } catch (e) {
    alert("Erro: " + (e.response?.data?.detail || e.message));
  }
};

const salvarItem = async () => {
  try {
    await api.createProduto(formItem);
    Object.assign(formItem, { nome: '', categoria: '', marca: '' });
    await carregarDados();
  } catch (e) {
    alert("Erro ao criar item.");
  }
};

const deletar = async (contexto, id) => {
  if(!confirm("Remover este registro?")) return;
  try {
    if(contexto === 'func') await api.deleteFuncionario(id);
    if(contexto === 'item') await api.deleteByEndpoint(`/produtos/${id}`);
    await carregarDados();
  } catch (e) {
    alert("Erro ao deletar.");
  }
};

onMounted(() => {
  carregarDados();
});
</script>

<template>
  <div class="nav-space"></div>
  <div class="page-container">

    <section class="card-panel">
      <div class="panel-header">
        <h2>Funcionários</h2>
      </div>

      <div class="form-row">
        <input v-model="formFuncionario.nome" placeholder="Nome Completo" />
        <input v-model="formFuncionario.CPF" placeholder="CPF (somente núm.)" maxlength="11" />

        <select v-model="formFuncionario.tipo">
          <option value="garcom">Garçom</option>
          <option value="seguranca">Segurança</option>
          <option value="caixa">Caixa</option>
          <option value="faxineiro">Faxineiro</option>
          <option value="balconista">Balconista</option>
        </select>

        <select v-model="formFuncionario.expediente">
          <option value="noite">Noite</option>
          <option value="manha">Manhã</option>
          <option value="tarde">Tarde</option>
        </select>

        <input type="number" v-model="formFuncionario.salario" placeholder="Salário R$" class="input-short" />
        <input type="date" v-model="formFuncionario.data_contratacao" />

        <button @click="salvarFuncionario" class="btn-action">+</button>
      </div>

      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>Nome</th>
              <th>CPF</th>
              <th>Função</th>
              <th>Expediente</th>
              <th>Salário</th>
              <th>Ações</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="func in funcionarios" :key="func.id">
              <td>{{ func.nome }}</td>
              <td>{{ func.CPF }}</td>
              <td><span class="badge">{{ func.tipo }}</span></td>
              <td>{{ func.expediente }}</td>
              <td class="money">R$ {{ func.salario.toFixed(2) }}</td>
              <td>
                <button @click="deletar('func', func.id)" class="btn-delete">Excluir</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <section class="card-panel mt-large">
      <div class="panel-header">
        <h2>Itens & Estrutura</h2>
        <p class="subtitle">Gerenciamento de insumos, mobília e descartáveis</p>
      </div>

      <div class="form-row">
        <input v-model="formItem.nome" placeholder="Nome do Item (ex: Copo Descartável)" />
        <input v-model="formItem.categoria" placeholder="Categoria (ex: Descartáveis)" />
        <input v-model="formItem.marca" placeholder="Marca" />
        <button @click="salvarItem" class="btn-action">+</button>
      </div>

      <div class="grid-items">
        <div v-for="item in itensEstrutura" :key="item.id" class="item-card">
          <div class="item-info">
            <h3>{{ item.nome }}</h3>
            <p>{{ item.categoria }} | {{ item.marca }}</p>
          </div>
          <div class="item-qtd">
            <span class="qtd-label">Qtd:</span>
            <span class="qtd-value">{{ item.quantidade }}</span>
          </div>
          <button @click="deletar('item', item.id)" class="btn-close">×</button>
        </div>
      </div>
    </section>

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

/* --- LAYOUT DA PÁGINA --- */
.nav-space {
    background-image: linear-gradient(
        220deg,
        var(--edna-green),
        var(--edna-blue)
    );
}

.page-container {
    background-color: var(--edna-black);
    min-height: 100vh;
    padding: 40px;
    display: flex;
    flex-direction: column;
    gap: 40px;
    font-family: "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
    color: var(--edna-white);
}

/* --- PAINÉIS (CARDS GRANDES) --- */
.card-panel {
    background-color: var(--edna-dark-gray);
    border: 1px solid var(--edna-gray);
    border-radius: 12px;
    padding: 25px;
    box-shadow: 0 4px 10px rgba(0,0,0,0.3);
}

.mt-large {
    margin-top: 20px;
}

.panel-header {
    margin-bottom: 25px;
    border-bottom: 1px solid var(--edna-gray);
    padding-bottom: 10px;
}

.panel-header h2 {
  margin-top: 0;
  color: var(--edna-yellow);
  font-size: 1.3rem;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.subtitle {
    color: var(--edna-light-gray);
    margin-top: 5px;
}

/* --- FORMULÁRIOS --- */
.form-row {
    display: flex;
    gap: 15px;
    flex-wrap: wrap;
    background-color: var(--edna-dark-gray); /* Fundo mais escuro para destaque */
    padding: 20px 0px;
    border-radius: 8px;
    margin-bottom: 25px;
    align-items: center;
}

input, select {
    background-color: var(--edna-gray);
    color: var(--edna-white);
    border: 1px solid var(--edna-gray);
    padding: 12px;
    border-radius: 6px;
    flex: 1;
    min-width: 120px;
    outline: none;
    font-size: 0.95rem;
    box-sizing: border-box;
}

input:focus,
select:focus {
    border: 2px solid var(--edna-orange);
}

.input-short {
    max-width: 120px;
}

/* Botão de Adicionar (+) */
.btn-action {
    background-color: var(--edna-green);
    color: var(--edna-black);
    font-weight: bold;
    width: 2.2rem;
    height: 2.2rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 8px;
    border: none;
    cursor: pointer;
    transition: filter 0.2s;
    flex-shrink: 0;
}

.btn-action:hover {
    filter: brightness(1.1);
}

/* --- TABELA DE FUNCIONÁRIOS --- */
.table-wrapper {
    overflow-x: auto;
    border-radius: 8px;
    border: 1px solid var(--edna-gray);
}

table {
    width: 100%;
    border-collapse: collapse;
    background-color: var(--edna-dark-gray);
}

th {
    background-color: var(--edna-gray);
    color: var(--edna-light-gray);
    text-align: left;
    padding: 15px;
    font-weight: 600;
    text-transform: uppercase;
    font-size: 0.85rem;
}

td {
    border-bottom: 1px solid var(--edna-gray);
    padding: 15px;
    color: var(--edna-white);
}

tr:last-child td {
    border-bottom: none;
}

tr:hover td {
    background-color: rgba(200, 200, 255, 0.03);
}

/* Badges e Tags */
.badge {
    background-color: var(--edna-yellow);
    color: var(--edna-dark-gray);
    padding: 4px 10px;
    border-radius: 4px;
    font-size: 0.8rem;
    text-transform: uppercase;
    font-weight: bold;
    border: 1px solid var(--edna-gray);
}

.money {
    color: var(--edna-green);
    font-family: monospace;
    font-size: 1rem;
}

/* Botões de Ação da Tabela */
.btn-delete {
    background-color: transparent;
    border: 1px solid var(--edna-red);
    color: var(--edna-red);
    padding: 6px 12px;
    border-radius: 4px;
    font-size: 0.8rem;
    cursor: pointer;
    transition: all 0.2s;
    text-transform: uppercase;
    font-weight: bold;
}

.btn-delete:hover {
    background-color: var(--edna-red);
    color: var(--edna-dark-gray);
}

/* --- GRID DE ITENS (Cards) --- */
.grid-items {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(312px, 1fr));
    gap: 20px;
}

.item-card {
    background-color: var(--edna-gray);
    border-radius: 8px;
    padding: 20px;
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    position: relative;
    transition: transform 0.2s;
    border-left: 4px solid var(--edna-green);
}

.item-card:hover {
    border-color: var(--edna-white);
    transform: translateY(-3px);
}

.item-info h3 {
    margin: 0;
    font-size: 1.1rem;
    color: var(--edna-white);
    margin-bottom: 5px;
}

.item-info p {
    margin: 0;
    font-size: 0.85rem;
    color: var(--edna-light-gray);
    text-transform: uppercase;
}

.item-qtd {
    text-align: center;
    background-color: var(--edna-dark-gray);
    padding: 8px 12px;
    border-radius: 6px;
    min-width: 60px;
    margin-top: 1rem;
    margin-left: 1rem;
}

.qtd-label {
    display: block;
    font-size: 0.7rem;
    color: var(--edna-light-gray);
    text-transform: uppercase;
}

.qtd-value {
    font-size: 1rem;
    color: var(--edna-white);
    font-weight: bold;
}

.btn-close {
    position: absolute;
    top: 8px;
    right: 8px;
    background: none;
    border: none;
    color: var(--edna-red);
    opacity: 0.7;
    font-size: 1.5rem;
    line-height: 1;
    cursor: pointer;
    transition: color 0.2s;
}

.btn-close:hover {
    color: var(--edna-red);
    transform: scale(1.2);
    opacity: 1.0;
}

/* Scrollbar da página */
::-webkit-scrollbar { width: 8px; }
::-webkit-scrollbar-track { background: var(--edna-black); }
::-webkit-scrollbar-thumb { background: var(--edna-gray); border-radius: 4px; }
::-webkit-scrollbar-thumb:hover { background: var(--edna-light-gray); }
</style>
