-- Rollback para 000011_insert_impressoes.up.sql
-- Remover impressões inseridas
DELETE FROM imprime
WHERE lisbn IN (
  '978-3-16-148410-0',
  '978-3-16-148411-7',
  '978-3-16-148412-4',
  '978-3-16-148413-1',
  '978-3-16-148414-8',
  '978-3-16-148415-5',
  '978-3-16-148416-2',
  '978-3-16-148417-9',
  '978-3-16-148418-6',
  '978-3-16-148419-3'
;

-- Remover contratos relacionados às gráficas contratadas
DELETE FROM Contrato
WHERE grafica_cont_id IN (2, 3)
  AND nome_responsavel IN ('Carlos Almeida', 'Maria Ferreira')
;

-- Remover entradas de Contratada
DELETE FROM Contratada
WHERE grafica_id IN (2, 3)
  AND endereco IN ('Avenida Central, 111', 'Rua dos Cravos, 222')
;

-- Remover entrada de Particular
DELETE FROM Particular
WHERE grafica_id = 1
;

-- Remover gráficas inseridas
DELETE FROM Grafica
WHERE nome IN ('Grafica Alfa', 'Grafica Beta', 'Grafica Gamma');
