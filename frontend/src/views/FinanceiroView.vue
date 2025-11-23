<script setup>
import { ref, reactive, computed, onMounted } from 'vue';
import api from '@/services/api';

// --- ESTADO ---
const loading = ref(false);
const periodo = reactive({
  start: new Date(new Date().getFullYear(), 0, 1).toISOString().split('T')[0], 
  end: new Date().toISOString().split('T')[0]
});

// Dados crus da API
const rawFinanceiro = ref(null);
const rawFolha = ref(null);
const rawLotes = ref([]);

// Caches para tradução de IDs (Novo)
const produtosMap = ref({});
const fornecedoresMap = ref({});

// Controle de UI
const mesesExpandidos = ref({});

// --- AÇÕES ---

const carregarDados = async () => {
  if (!periodo.start || !periodo.end) return alert("Selecione o período.");
  
  loading.value = true;
  try {
    // 1. Buscamos os relatórios e TAMBÉM os dados auxiliares (Produtos e Fornecedores)
    const [resFin, resFolha, resLotes, resProds, resForns] = await Promise.all([
      api.getFinancialReport({ start: periodo.start, end: periodo.end, granularity: 'month' }),
      api.getPayrollReport({ start: periodo.start, end: periodo.end }),
      api.getLotes(), // Traz todos para filtrar depois
      api.getProdutos(), // Traz lista geral de produtos para pegar o nome
      api.getFornecedores() // Traz lista de fornecedores
    ]);

    rawFinanceiro.value = resFin.data;
    rawFolha.value = resFolha.data;
    rawLotes.value = resLotes.data || [];

    // 2. Criar Mapas de Hash para acesso rápido por ID
    // Ex: { 1: "Cerveja", 2: "Fritas" }
    const pMap = {};
    (resProds.data || []).forEach(p => pMap[p.id] = p.nome);
    produtosMap.value = pMap;

    const fMap = {};
    (resForns.data || []).forEach(f => fMap[f.id] = f.nome);
    fornecedoresMap.value = fMap;

  } catch (error) {
    console.error("Erro ao carregar relatórios:", error);
    alert("Erro ao carregar dados.");
  } finally {
    loading.value = false;
  }
};

// --- PROCESSAMENTO (COMPUTED) ---
const relatorioUnificado = computed(() => {
  if (!rawFinanceiro.value || !rawFinanceiro.value.series) return [];

  return rawFinanceiro.value.series.map(ponto => {
    const mesChave = ponto.date.substring(0, 7); // "YYYY-MM"

    // 1. Filtrar Lotes deste mês
    const lotesDoMes = rawLotes.value.filter(lote => {
      if (!lote.data_fornecimento) return false;
      return lote.data_fornecimento.substring(0, 7) === mesChave;
    }).map(l => {
      const nomeProduto = produtosMap.value[l.id_produto] || `Produto ID ${l.id_produto}`;
      const nomeFornecedor = fornecedoresMap.value[l.id_fornecedor] || `Fornecedor ID ${l.id_fornecedor}`;

      return {
        id: l.id_lote,
        nomeProduto: nomeProduto,
        nomeFornecedor: nomeFornecedor,
        descricao: `${nomeProduto}`, 
        subDescricao: `Fornecedor: ${nomeFornecedor}`,
        valor: l.preco_unitario * l.quantidade_inicial,
        data: l.data_fornecimento
      };
    });

    // 2. Encontrar Folha deste mês
    const folhaDoMes = rawFolha.value?.folhas_por_mes?.find(f => {
      const dateObj = new Date(mesChave + "-02");
      const mesNome = dateObj.toLocaleString('en-US', {month:'long'}); 
      return f.ano === dateObj.getFullYear() && f.mes.toLowerCase().includes(mesNome.toLowerCase());
    });

    // --- FIX: FILTRO ROBUSTO DE DATA ---
    const salariosDoMes = folhaDoMes ? folhaDoMes.funcionarios
      .filter(func => {
        // Se não tiver data, não exibe (segurança contra dados corrompidos)
        if (!func.data_contratacao) return false;

        // Cria datas reais para comparação
        // Adicionamos "T00:00:00" para garantir que o navegador interprete como data local/ISO corretamente sem subtrair fuso
        const dtContratacao = new Date(func.data_contratacao);
        const dtRelatorio = new Date(mesChave + "-01T00:00:00");

        // Compara apenas Ano e Mês
        const anoCont = dtContratacao.getFullYear();
        const mesCont = dtContratacao.getMonth(); // 0 a 11

        const anoRel = dtRelatorio.getFullYear();
        const mesRel = dtRelatorio.getMonth();    // 0 a 11

        // Lógica: O mês do relatório deve ser IGUAL ou POSTERIOR ao mês de contratação
        if (anoRel > anoCont) return true;
        if (anoRel === anoCont && mesRel >= mesCont) return true;
        
        return false;
      })
      .map(func => ({
        id: func.id_funcionario,
        nome: func.nome,
        cargo: func.tipo,
        valor: func.salario_total
      })) : [];

    // 3. Calcular Totais (Corrigindo o Bug do 0)
    const totalLotes = lotesDoMes.reduce((acc, l) => acc + l.valor, 0);
    const totalSalarios = salariosDoMes.reduce((acc, s) => acc + s.valor, 0);
    
    const despesaReal = totalLotes + totalSalarios;
    const saldoReal = ponto.receita - despesaReal;

    return {
      mesAno: mesChave,
      receitaTotal: ponto.receita,
      despesaTotal: despesaReal, 
      saldo: saldoReal,
      detalhes: {
        lotes: lotesDoMes,
        salarios: salariosDoMes,
        somaLotes: totalLotes,
        somaSalarios: totalSalarios
      }
    };
  }).reverse();
});

// --- HELPERS ---
const toggleDetalhes = (mesAno) => {
  mesesExpandidos.value[mesAno] = !mesesExpandidos.value[mesAno];
};

const formatMoney = (val) => {
  return val.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' });
};

const formatMes = (isoMes) => {
  const [ano, mes] = isoMes.split('-');
  const date = new Date(parseInt(ano), parseInt(mes) - 1, 2);
  // Ajuste para pt-BR ficar bonito (ex: NOVEMBRO DE 2025)
  const mesNome = date.toLocaleDateString('pt-BR', { month: 'long' }).toUpperCase();
  const anoNum = date.getFullYear();
  return { mes: mesNome, ano: anoNum };
};

onMounted(() => {
  carregarDados();
});
</script>

<template>
  <div class="nav-space"></div>
  <div class="page-container">
    
    <section class="card-panel header-panel">
      <div class="panel-header">
        <h2>Relatórios Financeiros</h2>
        <p class="subtitle">Visão consolidada de Vendas, Compras e Pagamentos</p>
      </div>

      <div class="filters">
        <div class="date-group">
          <label>Início</label>
          <input type="date" v-model="periodo.start">
        </div>
        <div class="date-group">
          <label>Fim</label>
          <input type="date" v-model="periodo.end">
        </div>
        <button @click="carregarDados" class="btn-generate">
          <span v-if="!loading">Gerar Relatório</span>
          <span v-else>Carregando...</span>
        </button>
      </div>
    </section>

    <div class="reports-list">
      <div v-if="relatorioUnificado.length === 0" class="empty-state">
        Nenhum dado encontrado para o período.
      </div>

      <div 
        v-for="item in relatorioUnificado" 
        :key="item.mesAno" 
        class="month-card"
        :class="{ expanded: mesesExpandidos[item.mesAno] }"
      >
        <div class="month-header" @click="toggleDetalhes(item.mesAno)">
          <div class="month-title">
            <div class="calendar-icon-box">
              <span class="cal-month">📅</span>
            </div>
            <div class="title-text">
              <h3>{{ formatMes(item.mesAno).mes }}</h3>
              <span class="title-year">DE {{ formatMes(item.mesAno).ano }}</span>
            </div>
          </div>

          <div class="summary-grid">
            <div class="summary-item">
              <span class="label">Entradas (Vendas)</span>
              <span class="value green">{{ formatMoney(item.receitaTotal) }}</span>
            </div>
            <div class="summary-item">
              <span class="label">Saídas (Gastos)</span>
              <span class="value red">{{ formatMoney(item.despesaTotal) }}</span>
            </div>
            <div class="summary-item highlight">
              <span class="label">Saldo Final</span>
              <span class="value" :class="item.saldo >= 0 ? 'green' : 'red'">
                {{ formatMoney(item.saldo) }}
              </span>
            </div>
          </div>

          <button class="btn-toggle">
            {{ mesesExpandidos[item.mesAno] ? '▲' : '▼' }}
          </button>
        </div>

        <div v-if="mesesExpandidos[item.mesAno]" class="month-details">
          
          <div class="details-section">
            <h4>Detalhamento de Gastos</h4>
            
            <div class="details-columns">
              <div class="detail-column">
                <div class="col-header">
                  <h5>Compras de Estoque (Lotes)</h5>
                  <span class="subtotal">{{ formatMoney(item.detalhes.somaLotes) }}</span>
                </div>
                <ul class="detail-list">
                  <li v-for="lote in item.detalhes.lotes" :key="lote.id">
                    <div class="li-info">
                      <span class="li-name">{{ lote.nomeProduto }}</span>
                      <span class="li-tag">{{ lote.nomeFornecedor }}</span>
                      <span class="li-date">{{ new Date(lote.data).toLocaleDateString() }}</span>
                    </div>
                    <span class="li-value red">- {{ formatMoney(lote.valor) }}</span>
                  </li>
                  <li v-if="item.detalhes.lotes.length === 0" class="empty-li">Nenhuma compra registrada.</li>
                </ul>
              </div>

              <div class="detail-column">
                <div class="col-header">
                  <h5>Folha de Pagamento</h5>
                  <span class="subtotal">{{ formatMoney(item.detalhes.somaSalarios) }}</span>
                </div>
                <ul class="detail-list">
                  <li v-for="func in item.detalhes.salarios" :key="func.id">
                    <div class="li-info">
                      <span class="li-name">{{ func.nome }}</span>
                      <span class="li-tag">{{ func.cargo }}</span>
                    </div>
                    <span class="li-value red">- {{ formatMoney(func.valor) }}</span>
                  </li>
                  <li v-if="item.detalhes.salarios.length === 0" class="empty-li">Nenhum pagamento registrado.</li>
                </ul>
              </div>
            </div>
          </div>

          <div class="details-section revenue-section">
            <h4>Detalhamento de Ganhos</h4>
            <div class="revenue-block">
              <p>Total bruto de vendas realizadas no período:</p>
              <span class="big-revenue">{{ formatMoney(item.receitaTotal) }}</span>
            </div>
          </div>

        </div>
      </div>
    </div>

  </div>
</template>

<style scoped>
/* Mantenha o mesmo CSS da resposta anterior, ele está correto para o layout */
/* --- PALETA E.D.N.A. --- */
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

/* --- LAYOUT --- */
.nav-space {
    background-image: linear-gradient(220deg, var(--edna-blue), var(--edna-green));
}

.page-container {
    background-color: var(--edna-black);
    min-height: 100vh;
    padding: 40px;
    color: var(--edna-white);
    display: flex;
    flex-direction: column;
    gap: 30px;
}

/* --- HEADER & FILTERS --- */
.card-panel {
    background-color: var(--edna-dark-gray);
    border: 1px solid var(--edna-gray);
    border-radius: 12px;
    padding: 25px;
    box-shadow: 0 4px 10px rgba(0,0,0,0.3);
}

.panel-header h2 {
    color: var(--edna-yellow);
    font-size: 1.5rem;
    margin: 0 0 5px 0;
    font-weight: 300;
    text-transform: uppercase;
    letter-spacing: 1px;
}

.subtitle {
    color: var(--edna-light-gray);
    font-size: 1rem;
    margin-bottom: 20px;
}

.filters {
    display: flex;
    gap: 20px;
    align-items: flex-end;
    background-color: var(--edna-gray);
    padding: 15px;
    border-radius: 8px;
}

.date-group {
    display: flex;
    flex-direction: column;
}

.date-group label {
    font-size: 0.8rem;
    color: var(--edna-light-gray);
    margin-bottom: 5px;
}

input[type="date"] {
    background-color: var(--edna-gray);
    border: 2px solid var(--edna-dark-gray);
    color: var(--edna-white);
    padding: 10px;
    border-radius: 6px;
    outline: none;
    font-family: inherit;
}

input[type="date"]:focus {
    border: 2px solid var(--edna-orange);
    border-color: var(--edna-orange);
}

.btn-generate {
    background-color: var(--edna-green);
    border: none;
    padding: 10px 25px;
    border-radius: 6px;
    font-weight: bold;
    cursor: pointer;
    height: 2.2rem;
    transition: filter 0.2s;
}

.btn-generate span {
  color: var(--edna-dark-gray);
}

.btn-generate:hover {
    filter: brightness(1.1);
}

/* --- MONTH CARDS --- */
.reports-list {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.month-card {
    background-color: var(--edna-dark-gray);
    border: 1px solid var(--edna-gray);
    border-radius: 10px;
    overflow: hidden;
    transition: all 0.3s;
}

.month-card:hover {
    border-color: var(--edna-light-gray);
}

.month-card.expanded {
    border: 2px solid var(--edna-blue);
    box-shadow: 0 0 15px rgba(0,0,0,0.2);
}

.month-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 30px;
    cursor: pointer;
}

.month-title {
    display: flex;
    align-items: center;
    gap: 15px;
    width: 250px;
}

.calendar-icon-box {
    padding: 5px;
    border-radius: 8px;
}

.title-text h3 {
    margin: 0;
    font-size: 1.4rem;
    color: var(--edna-white);
    line-height: 1.1;
}

.title-year {
    font-size: 1rem;
    color: var(--edna-light-gray);
}

.summary-grid {
    display: flex;
    gap: 40px;
    flex: 1;
    justify-content: center;
}

.summary-item {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.summary-item .label {
    font-size: 0.8rem;
    color: var(--edna-light-gray);
    text-transform: uppercase;
    margin-bottom: 0.3rem;
}

.summary-item .value {
    font-size: 1.2rem;
}

.summary-item.highlight {
    border-left: 1px solid var(--edna-gray);
    padding-left: 40px;
}

.btn-toggle {
    background: none;
    border: none;
    color: var(--edna-light-gray);
    font-size: 1rem;
    cursor: pointer;
}

/* --- DETALHES (Area Expandida) --- */
.month-details {
    background-color: var(--edna-black);
    padding: 30px;
    animation: slideDown 0.3s ease-out;
}

@keyframes slideDown {
    from { opacity: 0; transform: translateY(-10px); }
    to { opacity: 1; transform: translateY(0); }
}

.details-section h4 {
    color: var(--edna-yellow);
    margin-top: 0;
    padding-bottom: 10px;
    margin-bottom: 20px;
    font-size: 1.2rem;
}

.details-columns {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 30px;
    margin-bottom: 30px;
}

.detail-column {
    background-color: var(--edna-dark-gray);
    border-radius: 8px;
    padding: 15px;
}

.col-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;
}

.col-header h5 {
    margin: 0;
    font-size: 1.1rem;
    color: var(--edna-light-gray);
}

.subtotal {
    background-color: var(--edna-black);
    padding: 8px 10px;
    border-radius: 4px;
    font-size: 0.9rem;
    color: var(--edna-red);
}

/* Listas internas */
.detail-list {
    list-style: none;
    padding: 0;
    margin: 0;
    max-height: 250px;
    overflow-y: auto;
}

.detail-list li {
    display: flex;
    justify-content: space-between;
    padding: 10px 0;
    margin: 0rem 1rem;
    border-bottom: 2px dashed var(--edna-light-gray);
    font-size: 0.9rem;
}

.detail-list li:last-child { border-bottom: none; }

.li-info { display: flex; flex-direction: column; }
.li-name { color: var(--edna-white); font-weight: bold; }
.li-tag { font-size: 0.8rem; color: var(--edna-light-gray); margin-top: 2px; }
.li-date { font-size: 0.8rem; color: var(--edna-light-gray); margin-top: 2px; }

.empty-li {
    font-style: italic;
    color: var(--edna-gray);
    text-align: center;
    padding: 10px;
}

/* Seção de Receita */
.revenue-section {
    background-color: rgba(90, 211, 176, 0.1); /* Green tint */
    padding: 20px;
    border-radius: 8px;
    border: 1px solid rgba(90, 211, 176, 0.2);
}

.revenue-block {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.revenue-block p { margin: 0; color: var(--edna-white); }

.big-revenue {
    font-size: 1.5rem;
    color: var(--edna-green);
    font-weight: bold;
}

/* Helpers de Cor */
.green { color: var(--edna-green); }
.red { color: var(--edna-red); }

/* Scrollbar Interna */
.detail-list::-webkit-scrollbar { width: 6px; }
.detail-list::-webkit-scrollbar-thumb { background: var(--edna-gray); border-radius: 3px; }
.detail-list::-webkit-scrollbar-track { background: var(--edna-black); }

.empty-state {
    text-align: center;
    color: var(--edna-light-gray);
    margin-top: 40px;
    font-style: italic;
}
</style>
