// src/services/api.js
import axios from "axios";

// Cria uma instância do axios que já aponta para o seu back-end
const apiClient = axios.create({
  baseURL: "http://localhost:8080/api/v1", // O BasePath que definimos!
  headers: {
    "Content-Type": "application/json",
  },
});

export default {
  getFornecedores(filters = null) {
    return apiClient.get("/fornecedores");
  },
  getFornecedorById(id) {
    return apiClient.get(`/fornecedores/${id}`);
  },
  createFornecedor(data) {
    // data é um objeto JS, ex: { nome: "...", cnpj: "..." }
    return apiClient.post("/fornecedores", data);
  },
  deleteFornecedor(id) {
    return apiClient.delete(`/fornecedores/${id}`);
  },
  getProdutos(filters = null) {
    // filters pode ser um objeto, ex: { params: { 'filter-nome': 'ilike.Cerveja' } }
    return apiClient.get("/produtos", filters);
  },
  getProdutosComerciais(filters = null) {
    // filters pode ser um objeto, ex: { params: { 'filter-nome': 'ilike.Cerveja' } }
    return apiClient.get("/produtos/comercial", filters);
  },
  getOfertas(filters = null) {
    return apiClient.get("/ofertas");
  },
  createProdutoComercial(data) {
    return apiClient.post("/produtos/comercial", data);
  },
  getProdutoQtd(id) {
    return apiClient.get(`/produtos/quantidade/${id}`);
  },
  createOferta(data) {
    return apiClient.post("/ofertas", data);
  },
  deleteByEndpoint(endpoint) {
    return apiClient.delete(endpoint);
  },
  // --- FUNCIONÁRIOS ---
  getFuncionarios(filters = null) {
    return apiClient.get("/funcionarios", filters);
  },
  createFuncionario(data) {
    return apiClient.post("/funcionarios", data);
  },
  deleteFuncionario(id) {
    return apiClient.delete(`/funcionarios/${id}`);
  },
  // --- PRODUTOS ESTRUTURAIS ---
  getProdutosEstruturais(filters = null) {
    return apiClient.get("/produtos/estrutural", filters);
  },
  // A criação de produto é genérica (estrutural ou comercial)
  createProduto(data) {
    return apiClient.post("/produtos", data);
  },
  createVenda(data) {
    // data: { id_cliente, id_funcionario, tipo_pagamento, ... }
    return apiClient.post("/vendas", data);
  },
  // O backend espera que criemos os itens associados à venda
  createItemVenda(data) {
    // data: { id_venda, id_produto, id_lote (opcional dependendo da lógica), quantidade, valor_unitario }
    return apiClient.post("/item_venda", data);
  },
  // Auxiliar para buscar lote disponível (Lógica FIFO necessária para o item_venda)
  // Nota: Se o seu backend não faz isso automático, o front precisa descobrir o lote.
  // Vou assumir por enquanto que vamos listar lotes ou pegar o primeiro disponível.
  getLotesPorProduto(idProduto) {
    return apiClient.get(`/lotes/produtos/${idProduto}`);
  },
  // --- VENDAS & HISTÓRICO ---
  getVendas(filters = null) {
    return apiClient.get("/vendas", filters);
  },

  // Para buscar os itens de uma venda específica
  getItemVenda(filters = null) {
    // Exemplo de uso: { params: { 'filter-id_venda': 'eq.1' } }
    return apiClient.get("/item_venda", filters);
  },
  // Necessário para descobrir o nome do produto a partir do item vendido (que só tem id_lote)
  getLote(id) {
    return apiClient.get(`/lotes/${id}`);
  },

  getLotes(filters = null) {
    return apiClient.get("/lotes", { params: filters });
  },
  createLote(data) {
    // data espera: { id_fornecedor, id_produto, data_fornecimento, validade, preco_unitario, quantidade_inicial, estragados }
    return apiClient.post("/lotes", data);
  },
  deleteLote(id) {
    return apiClient.delete(`/lotes/${id}`);
  },
  getClienteSaldo(id) {
    return apiClient.get(`/clientes/${id}/saldo`);
  },
  getClientes(filters = null) {
    return apiClient.get("/clientes");
  },
  createCliente(data) {
    return apiClient.post("/clientes", data);
  },
  // --- NOVOS MÉTODOS PARA EDIÇÃO E REMOÇÃO DE CLIENTES ---
  updateCliente(id, data) {
    return apiClient.put(`/clientes/${id}`, data);
  },
  deleteCliente(id) {
    return apiClient.delete(`/clientes/${id}`);
  },

  getFinancialReport(params) {
    return apiClient.get("/relatorios/financeiro", { params });
  },
  getPayrollReport(params) {
    return apiClient.get("/relatorios/folha-pagamento", { params });
  },
  updateVenda(id, data) {
    // data deve conter todos os campos: id_cliente, id_funcionario, datas, etc.
    return apiClient.put(`/vendas/${id}`, data);
  },
};
