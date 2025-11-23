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
    return apiClient.post("/fornecedores", data);
  },
  deleteFornecedor(id) {
    return apiClient.delete(`/fornecedores/${id}`);
  },
  getProdutos(filters = null) {
    return apiClient.get("/produtos", filters);
  },
  getProdutosComerciais(filters = null) {
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
  getFuncionarios(filters = null) {
    return apiClient.get("/funcionarios", filters);
  },
  createFuncionario(data) {
    return apiClient.post("/funcionarios", data);
  },
  deleteFuncionario(id) {
    return apiClient.delete(`/funcionarios/${id}`);
  },
  getProdutosEstruturais(filters = null) {
    return apiClient.get("/produtos/estrutural", filters);
  },
  createProduto(data) {
    return apiClient.post("/produtos", data);
  },
  createVenda(data) {
    return apiClient.post("/vendas", data);
  },
  createItemVenda(data) {
    return apiClient.post("/item_venda", data);
  },
  getLotesPorProduto(idProduto) {
    return apiClient.get(`/lotes/produtos/${idProduto}`);
  },
  getVendas(filters = null) {
    return apiClient.get("/vendas", filters);
  },
  getItemVenda(filters = null) {
    return apiClient.get("/item_venda", filters);
  },
  getLote(id) {
    return apiClient.get(`/lotes/${id}`);
  },
  createOferta(data) {
    return apiClient.post("/ofertas", data);
  },
  deleteByEndpoint(endpoint) {
    return apiClient.delete(endpoint);
  },

  getLotes(filters = null) {
    return apiClient.get("/lotes", { params: filters });
  },
  createLote(data) {
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
  getClienteSaldo(id) {
    return apiClient.get(`/clientes/${id}/saldo`);
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
  updateProdutoComercial(id, data) {
    return apiClient.put(`/produtos/comercial/${id}`, data);
  },
  updateFuncionario(id, data) {
    return apiClient.put(`/funcionarios/${id}`, data);
  },
  updateProduto(id, data) {
    return apiClient.put(`/produtos/${id}`, data);
  },
  updateLote(id, data) {
    return apiClient.put(`/lotes/${id}`, data);
  },
};
