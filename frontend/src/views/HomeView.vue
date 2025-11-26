<script setup>
import { ref } from "vue";

// Controle do Modal
const showModal = ref(false);
const selectedTopic = ref(null);

// Dados detalhados para cada seção
const guideDetails = {
    financeiro: {
        title: "Financeiro",
        icon: "💰",
        content: `
      O módulo financeiro é o coração da gestão do bar. Aqui você pode:
      <ul style="margin-top: 10px; margin-left: 20px; line-height: 1.6;">
        <li>Acompanhar o <strong>Fluxo de Caixa</strong> mensal com gráficos intuitivos.</li>
        <li>Visualizar o total de <strong>Entradas</strong> (Vendas) e <strong>Saídas</strong> (Despesas com estoque e folha de pagamento).</li>
        <li>Analisar o <strong>Lucro Líquido</strong> real do estabelecimento.</li>
        <li>Filtrar relatórios por períodos específicos para tomadas de decisão.</li>
      </ul>
    `,
    },
    produtos: {
        title: "Produtos",
        icon: "🍺",
        content: `
      Gerencie todo o seu catálogo de itens. O sistema divide produtos em duas categorias:
      <ul style="margin-top: 10px; margin-left: 20px; line-height: 1.6;">
        <li><strong>Comerciais:</strong> Itens destinados à venda direta (Cervejas, Porções, Drinks). Você define o preço de venda aqui.</li>
        <li><strong>Estruturais:</strong> Insumos e materiais de uso interno (Copos, Guardanapos, Produtos de Limpeza).</li>
      </ul>
      Você também pode criar e gerenciar <strong>Ofertas e Combos</strong> promocionais nesta seção.
    `,
    },
    vendas: {
        title: "Vendas (PDV)",
        icon: "🧾",
        content: `
      A frente de caixa (Ponto de Venda) foi desenhada para agilidade:
      <ul style="margin-top: 10px; margin-left: 20px; line-height: 1.6;">
        <li><strong>Nova Venda:</strong> Selecione o cliente, o funcionário atendente e adicione produtos ao carrinho.</li>
        <li><strong>Controle de Estoque:</strong> A venda baixa automaticamente o estoque do lote mais antigo (FIFO).</li>
        <li><strong>Histórico:</strong> Consulte todas as vendas realizadas, filtre por data e veja detalhes dos itens vendidos.</li>
      </ul>
    `,
    },
    fornecedores: {
        title: "Fornecedores & Estoque",
        icon: "🚚",
        content: `
      Mantenha o controle da cadeia de suprimentos:
      <ul style="margin-top: 10px; margin-left: 20px; line-height: 1.6;">
        <li>Cadastre empresas parceiras com CNPJ e dados de contato.</li>
        <li><strong>Entrada de Lotes:</strong> Ao comprar produtos, você registra um novo lote vinculado a um fornecedor.</li>
        <li>Controle a <strong>validade</strong> e o preço de custo de cada lote individualmente.</li>
        <li>Monitore perdas registrando itens estragados ou avariados.</li>
      </ul>
    `,
    },
    funcionarios: {
        title: "Funcionários",
        icon: "👥",
        content: `
      Gestão completa da equipe do bar:
      <ul style="margin-top: 10px; margin-left: 20px; line-height: 1.6;">
        <li>Cadastre garçons, caixas, seguranças e equipe de limpeza.</li>
        <li>Defina salários, turnos (expediente) e datas de contratação.</li>
        <li>Estes dados são cruzados automaticamente com o Financeiro para cálculo da folha de pagamento mensal.</li>
      </ul>
    `,
    },
    clientes: {
        title: "Clientes",
        icon: "🍻", // Ícone ajustado para diferenciar de funcionários
        content: `
      Fidelize seus clientes e gerencie o crédito da casa:
      <ul style="margin-top: 10px; margin-left: 20px; line-height: 1.6;">
        <li>Cadastre clientes regulares para agilizar o atendimento.</li>
        <li><strong>Sistema de Fiado:</strong> Permite realizar vendas sem pagamento imediato.</li>
        <li><strong>Extrato de Dívidas:</strong> Visualize o saldo devedor de cada cliente e realize a quitação total ou parcial das dívidas através do histórico.</li>
      </ul>
    `,
    },
};

const openGuide = (key) => {
    selectedTopic.value = guideDetails[key];
    showModal.value = true;
};
</script>

<template>
    <div class="home-container">
        <div class="banner-wrapper">
            <img
                id="home-banner"
                src="/home_banner.png"
                alt="Banner EDNA Bar"
            />
            <div class="scroll-indicator">
                <p>Role para saber mais</p>
                <span class="arrow">▼</span>
            </div>
        </div>

        <div class="content-body">
            <section class="info-section">
                <h1 class="main-title">E.D.N.A.</h1>
                <h2 class="subtitle">Ecossistema De Negociação Alcoólica</h2>
                <p class="description">
                    O E.D.N.A é um sistema integrado de gerenciamento para
                    bares, focado em otimizar o controle de estoque, fluxo de
                    caixa e o relacionamento com clientes e fornecedores.
                </p>
            </section>

            <hr class="divider" />

            <section class="tutorial-section">
                <h2>Como Utilizar o Sistema</h2>
                <p class="hint-text">Clique nos cartões para ver detalhes</p>

                <div class="cards-grid">
                    <div class="card-guide" @click="openGuide('financeiro')">
                        <h3>💰 Financeiro</h3>
                        <p>
                            Visualize relatórios de lucros, despesas e projeções
                            futuras.
                        </p>
                    </div>

                    <div class="card-guide" @click="openGuide('produtos')">
                        <h3>🍺 Produtos</h3>
                        <p>
                            Gerencie cardápio, itens comerciais e estruturais.
                        </p>
                    </div>

                    <div class="card-guide" @click="openGuide('vendas')">
                        <h3>🧾 Vendas</h3>
                        <p>PDV ágil, histórico de transações e comandas.</p>
                    </div>

                    <div class="card-guide" @click="openGuide('fornecedores')">
                        <h3>🚚 Fornecedores</h3>
                        <p>Entrada de notas, gestão de lotes e parceiros.</p>
                    </div>

                    <div class="card-guide" @click="openGuide('funcionarios')">
                        <h3>👥 Funcionários</h3>
                        <p>
                            Gestão de equipe, cargos, turnos e folha salarial.
                        </p>
                    </div>

                    <div class="card-guide" @click="openGuide('clientes')">
                        <h3>🍻 Clientes</h3>
                        <p>
                            Cadastro de clientes, sistema de fiado e quitação de
                            dívidas.
                        </p>
                    </div>
                </div>
            </section>

            <div class="inspiration-section">
                <div class="gallery-container">
                    <img src="/edna_bar1.jpeg" alt="Foto do Bar da Edna 1" class="gallery-img">
                    <img src="/edna_bar2.jpeg" alt="Foto do Bar da Edna 2" class="gallery-img">
                </div>
                <p class="inspiration-caption">
                    Este site foi inspirado no bar da edna, o melhor da Bahia <span class="heart">❤️</span>
                </p>
            </div>

            <footer class="home-footer">
                <p>
                    &copy; 2025 E.D.N.A Bar System - Projeto de Banco de Dados I
                </p>
                <a class="gh-link" href="https://github.com/caio-bernardo/edna-bar">Github</a>
            </footer>
        </div>

        <Transition name="fade">
            <div
                v-if="showModal"
                class="modal-overlay"
                @click.self="showModal = false"
            >
                <div class="modal-content">
                    <header class="modal-header">
                        <span class="modal-icon">{{ selectedTopic.icon }}</span>
                        <h3>{{ selectedTopic.title }}</h3>
                        <button class="close-btn" @click="showModal = false">
                            ×
                        </button>
                    </header>
                    <div
                        class="modal-body"
                        v-html="selectedTopic.content"
                    ></div>
                </div>
            </div>
        </Transition>
    </div>
</template>

<style scoped>
.home-container {
    display: flex;
    flex-direction: column;
    width: 100%;
    background-color: var(--edna-black);
}

.banner-wrapper {
    position: relative;
    width: 100%;
    height: 85vh;
    display: flex;
    justify-content: center;
    z-index: 1;
}

#home-banner {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-bottom-left-radius: 5vw;
    border-bottom-right-radius: 5vw;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
    z-index: 0;
}

.scroll-indicator {
    position: absolute;
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
    text-align: center;
    color: var(--edna-white);
    z-index: 2;
    animation: bounce 2s infinite;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.8);
}

.scroll-indicator p {
    font-size: 0.9rem;
    margin-bottom: 5px;
    text-transform: uppercase;
    letter-spacing: 2px;
}

.arrow {
    font-size: 1.5rem;
}

.content-body {
    background-color: var(--edna-black);
    color: var(--edna-white);
    padding: 60px 10%;
    min-height: 60vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-top: -20px;
    padding-top: 60px;
}

.info-section {
    text-align: center;
    max-width: 800px;
    margin-bottom: 40px;
}

.main-title {
    font-size: 4rem;
    color: var(--edna-yellow);
    margin-bottom: 0.5rem;
    font-weight: bold;
    letter-spacing: 5px;
}

.subtitle {
    font-size: 1.5rem;
    color: var(--edna-orange);
    margin-bottom: 20px;
    text-transform: uppercase;
}

.description {
    font-size: 1.1rem;
    line-height: 1.6;
    color: var(--edna-light-gray);
}

.divider {
    width: 100px;
    border: 2px solid var(--edna-green);
    margin: 40px 0;
    border-radius: 2px;
}

.tutorial-section {
    width: 100%;
    text-align: center;
}

.tutorial-section h2 {
    font-size: 2rem;
    color: var(--edna-white);
    margin-bottom: 10px;
    border-bottom: 1px dashed var(--edna-gray);
    display: inline-block;
    padding-bottom: 10px;
}

.hint-text {
    color: var(--edna-light-gray);
    font-size: 0.9rem;
    margin-bottom: 30px;
    font-style: italic;
}

.cards-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 30px;
    width: 100%;
}

.card-guide {
    background-color: var(--edna-dark-gray);
    padding: 30px;
    border-radius: 12px;
    border: 1px solid var(--edna-gray);
    transition:
        transform 0.2s,
        border-color 0.2s,
        background-color 0.2s;
    text-align: left;
    cursor: pointer;
}

.card-guide:hover {
    transform: translateY(-5px);
    border-color: var(--edna-yellow);
    background-color: #353540;
}

.card-guide h3 {
    color: var(--edna-green);
    margin-bottom: 15px;
    font-size: 1.4rem;
    display: flex;
    align-items: center;
    gap: 10px;
}

.card-guide p {
    color: var(--edna-light-gray);
    font-size: 1rem;
    line-height: 1.4;
}

.home-footer {
    margin-top: 80px;
    color: var(--edna-gray);
    font-size: 0.8rem;
    text-align: center;
}

.gh-link {
    color: var(--edna-blue);
    text-decoration: none;
}

.gh-link:hover {
    text-decoration: underline;
}

@keyframes bounce {
    0%,
    20%,
    50%,
    80%,
    100% {
        transform: translateX(-50%) translateY(0);
    }
    40% {
        transform: translateX(-50%) translateY(-10px);
    }
    60% {
        transform: translateX(-50%) translateY(-5px);
    }
}

/* --- ESTILOS DO MODAL --- */
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
    z-index: 1000;
    backdrop-filter: blur(3px);
}

.modal-content {
    background-color: var(--edna-dark-gray);
    width: 90%;
    max-width: 600px;
    border-radius: 16px;
    border: 1px solid var(--edna-yellow);
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
    overflow: hidden;
    animation: slideUp 0.3s ease-out;
}

.modal-header {
    background-color: var(--edna-black);
    padding: 20px;
    display: flex;
    align-items: center;
    gap: 15px;
    border-bottom: 1px solid var(--edna-gray);
}

.modal-icon {
    font-size: 2rem;
}

.modal-header h3 {
    margin: 0;
    color: var(--edna-yellow);
    flex-grow: 1;
    font-size: 1.5rem;
}

.close-btn {
    background: none;
    border: none;
    color: var(--edna-light-gray);
    font-size: 2rem;
    cursor: pointer;
    line-height: 1;
}

.close-btn:hover {
    color: var(--edna-red);
}

.modal-body {
    padding: 30px;
    color: var(--edna-white);
    font-size: 1.1rem;
}

@keyframes slideUp {
    from {
        transform: translateY(20px);
        opacity: 0;
    }
    to {
        transform: translateY(0);
        opacity: 1;
    }
}

.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}

.inspiration-section {
  margin-top: 60px;
  text-align: center;
  width: 100%;
}

.gallery-container {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 30px;
  margin-bottom: 25px;
}

.gallery-img {
  height: 500px; 
  width: auto;
  
  object-fit: cover;
  border-radius: 12px;
  box-shadow: 0 6px 12px rgba(0,0,0,0.15);
  transition: transform 0.3s ease;
}

.gallery-img:hover {
  transform: scale(1.01); /* Leve zoom ao passar o mouse */
}

/* Texto da legenda abaixo das fotos */
.inspiration-caption {
  font-size: 1.5rem;
  color: var(--edna-green);
  font-weight: 500;
}

.heart {
  color: var(--edna-red); /* Cor vermelha para o coração */
}

/* Responsividade */
@media (max-width: 768px) {
    .main-title {
        font-size: 2.5rem;
    }
    .subtitle {
        font-size: 1.1rem;
    }
    .banner-wrapper {
        height: 60vh;
    }
    .content-body {
        padding: 40px 5%;
    }
    .inspiration-section {
        margin-top: 40px;
    }

    .gallery-container {
        flex-direction: column;
        gap: 20px;
    }

    .gallery-img {
        height: auto;
        width: 90%;
        max-height: 60vh;
    }

    .inspiration-caption {
        font-size: 1.1rem;
        padding: 0 15px;
    }
}
</style>
