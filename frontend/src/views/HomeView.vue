<script setup>
// Nenhuma lógica de script necessária por enquanto
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
                    Desenvolvido para trazer modernidade e eficiência para o seu
                    estabelecimento.
                </p>
            </section>

            <hr class="divider" />

            <section class="tutorial-section">
                <h2>Como Utilizar o Sistema</h2>
                <div class="cards-grid">
                    <div class="card-guide">
                        <h3>💰 Financeiro</h3>
                        <p>
                            Visualize relatórios de lucros, despesas e projeções
                            futuras. Acompanhe o desempenho mensal do bar.
                        </p>
                    </div>

                    <div class="card-guide">
                        <h3>🍺 Produtos</h3>
                        <p>
                            Gerencie seu cardápio. Adicione produtos comerciais
                            (venda) e estruturais (insumos/mobília).
                        </p>
                    </div>

                    <div class="card-guide">
                        <h3>🧾 Vendas</h3>
                        <p>
                            Realize vendas no balcão (PDV), gerencie comandas e
                            visualize o histórico de transações.
                        </p>
                    </div>

                    <div class="card-guide">
                        <h3>🚚 Fornecedores</h3>
                        <p>
                            Cadastre seus parceiros de negócio e controle a
                            entrada de novos lotes de estoque.
                        </p>
                    </div>

                    <div class="card-guide">
                        <h3>👥 Funcionários</h3>
                        <p>
                            Gerencie sua equipe, turnos, salários e funções
                            dentro do estabelecimento.
                        </p>
                    </div>

                    <div class="card-guide">
                        <h3>👥 Clientes</h3>
                        <p>
                            Cadastre clientes fiéis, gerencie dívidas (fiado) e
                            histórico de consumo.
                        </p>
                    </div>
                </div>
            </section>

            <footer class="home-footer">
                <p>
                    &copy; 2025 E.D.N.A Bar System - Projeto de Banco de Dados I
                </p>
            </footer>
        </div>
    </div>
</template>

<style scoped>
/* Container Geral */
.home-container {
    display: flex;
    flex-direction: column;
    width: 100%;
    background-color: var(--edna-black);
}

/* --- Banner --- */
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

/* Indicador de Scroll */
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

/* --- Corpo do Conteúdo --- */
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

/* Seção Sobre */
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

/* Seção Tutorial */
.tutorial-section {
    width: 100%;
    text-align: center;
}

.tutorial-section h2 {
    font-size: 2rem;
    color: var(--edna-white);
    margin-bottom: 40px;
    border-bottom: 1px dashed var(--edna-gray);
    display: inline-block;
    padding-bottom: 10px;
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
        transform 0.3s,
        border-color 0.3s;
    text-align: left;
}

.card-guide:hover {
    transform: translateY(-5px);
    border-color: var(--edna-yellow);
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

/* Footer */
.home-footer {
    margin-top: 80px;
    color: var(--edna-gray);
    font-size: 0.8rem;
    text-align: center;
}

/* Animação da Seta */
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
}
</style>
