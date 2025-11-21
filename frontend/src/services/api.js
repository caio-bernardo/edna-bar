// src/services/api.js
import axios from 'axios';

// Cria uma instância do axios que já aponta para o seu back-end
const apiClient = axios.create({
    baseURL: 'http://localhost:8080/api/v1', // O BasePath que definimos!
    headers: {
        'Content-Type': 'application/json',
    },
});

export default {
    // TODO: Resto das funções de acesso à API 
    getFornecedores(filters = null) {
        return apiClient.get('/fornecedores');
    },
    getFornecedorById(id) {
        return apiClient.get(`/fornecedores/${id}`);
    },
    createFornecedor(data) {
        // data é um objeto JS, ex: { nome: "...", cnpj: "..." }
        return apiClient.post('/fornecedores', data);
    },
    deleteFornecedor(id) {
        return apiClient.delete(`/fornecedores/${id}`);
    },
    getProdutos(filters = null) {
        // filters pode ser um objeto, ex: { params: { 'filter-nome': 'ilike.Cerveja' } }
        return apiClient.get('/produtos', filters);
    },
    getProdutosComerciais(filters = null) {
        // filters pode ser um objeto, ex: { params: { 'filter-nome': 'ilike.Cerveja' } }
        return apiClient.get('/produtos/comercial', filters);
    },
    getClientes(filters = null) {
        return apiClient.get('/clientes');
    },
    getOfertas(filters = null) {
        return apiClient.get('/ofertas')
    },
    createProdutoComercial(data) {
        return apiClient.post('/produtos/comercial', data);
    },
    getProdutoQtd(id) {
        return apiClient.get(`/produtos/quantidade/${id}`);
    },
    createOferta(data) {
        return apiClient.post('/ofertas', data);
    },
    deleteByEndpoint(endpoint) {
        return apiClient.delete(endpoint);
    },
    getLotes(filters = null) {
        return apiClient.get('/lotes', { params: filters });
    },
    createLote(data) {
        // data espera: { id_fornecedor, id_produto, data_fornecimento, validade, preco_unitario, quantidade_inicial, estragados }
        return apiClient.post('/lotes', data);
    },
    deleteLote(id) {
        return apiClient.delete(`/lotes/${id}`);
    },
    // Auxiliar para popular o select de produtos no form de lotes
    getProdutosEstruturais(filters = null) {
        return apiClient.get('/produtos', { params: filters }); 
    }
};
