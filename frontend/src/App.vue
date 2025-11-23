<script setup></script>

<template>
    <div id="layout">
        <header>
            <nav class="navbar">
                <ul class="nav-links">
                    <li><RouterLink to="/">Início</RouterLink></li>
                    <li><RouterLink to="/financeiro">Financeiro</RouterLink></li>
                    <li><RouterLink to="/funcionarios">Funcionários</RouterLink></li>
                    <li><RouterLink to="/vendas">Vendas</RouterLink></li>
                    <li><RouterLink to="/produtos">Produtos</RouterLink></li>
                    <li><RouterLink to="/fornecedores">Fornecedores</RouterLink></li>
                </ul>
            </nav>
        </header>

        <main class="content">
            <RouterView />
        </main>
    </div>
</template>

<style>
.navbar {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    background-color: transparent;
    color: white;
    position: absolute;
    width: 100vw;
    border: 0px;
    margin: 1rem 0rem;
    z-index: 100; /* ADICIONADO: Garante que o menu fique sobre o banner da Home */
}
.nav-brand {
    font-weight: bold;
}
.nav-links {
    list-style: none;
    display: flex;
    padding: 0px 2rem;
}
.nav-links li {
    padding: 1vh 2vw;
}
.nav-links a {
    color: var(--edna-black);
    text-decoration: none;
    font-weight: 500;
}
/* Ajuste para garantir leitura em fundos escuros/banners se necessário,
   mas mantendo a lógica original */
.nav-links a:hover,
.nav-links a.router-link-active {
    color: var(--edna-yellow);
}
</style>
