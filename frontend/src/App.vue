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
    align-items: center;
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 10vh;
    
    margin: 0;
    border: 0;
    padding: 0 20px;
    box-sizing: border-box;
    
    z-index: 100;
    background-color: transparent;
    
    justify-content: flex-end; 
}

.nav-links {
    list-style: none;
    display: flex;
    padding: 0;
    margin: 0;
    gap: 15px;
    align-items: center;
    height: 100%;
}

.nav-links li {
    flex-shrink: 0;
}

.nav-links a {
    color: var(--edna-black);
    text-decoration: none;
    font-weight: 500;
    font-size: 1rem;
    white-space: nowrap; 
    transition: color 0.2s;
}

.nav-links a:hover,
.nav-links a.router-link-active {
    color: var(--edna-yellow);
    text-shadow: 0 1px 2px rgba(0,0,0,0.3);
}

@media (max-width: 768px) {
    .navbar {
        justify-content: flex-start;
        
        overflow-x: auto;
        overflow-y: hidden;
        
        -webkit-overflow-scrolling: touch; 
        
        scrollbar-width: none;
    }
    
    .navbar::-webkit-scrollbar {
        display: none;
    }

    .nav-links {
        flex-wrap: nowrap;
        padding-right: 20px;
    }
    
    .nav-links a {
        padding: 10px 5px;
    }
}
</style>
