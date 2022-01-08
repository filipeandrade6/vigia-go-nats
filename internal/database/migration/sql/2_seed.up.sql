INSERT INTO usuarios (usuario_id, email, funcao, senha) VALUES
    ('d03307d4-2b28-4c23-a004-3da25e5b8ce2', 'filipe@andrade.com.br', '{ADMIN, MANAGER, USER}', '$2a$10$n531epIH68yygcV6sNNqluZtyPc3smWxbw1WoWDDhOIqUP1Py/GTq'),
    ('d03307d4-2b28-4c23-a004-3da25e5b8cf3', 'filipe@andrade2.com.br', '{ADMIN, MANAGER, USER}', '$2a$10$n531epIH68yygcV6sNNqluZtyPc3smWxbw1WoWDDhOIqUP1Py/GTq')
    ON CONFLICT DO NOTHING;

INSERT INTO cameras (camera_id, descricao, endereco_ip, porta, canal, usuario, senha, latitude, longitude) VALUES
    ('d03307d4-2b28-4c23-a004-3da25e5b8ce3', 'desc 1', '10.20.30.40', '1', '1', 'admin', 'admin', '-12.2332', '-42.231'),
    ('d03307d4-2b28-4c23-a004-3da25e5b8aa3', 'desc 2', '45.56.78.89', '1', '1', 'admin', 'admin', '-12.2332', '-42.231')
    ON CONFLICT DO NOTHING;

INSERT INTO servidores_gravacao (servidor_gravacao_id, endereco_ip, porta, armazenamento, horas_retencao) VALUES
    ('d03307d4-2b28-4c23-a004-3da25e5b8bb1', '12.34.67.89', '6543', '/', '1'),
    ('d03307d4-2b28-4c23-a004-3da25e524bb1', '21.43.76.98', '3456', '/', '1')
    ON CONFLICT DO NOTHING;

INSERT INTO processos (processo_id, servidor_gravacao_id, camera_id, processador, adaptador) VALUES
    ('d03307d4-2b28-4c23-a004-3da32e5b8bb1', 'd03307d4-2b28-4c23-a004-3da25e5b8bb1', 'd03307d4-2b28-4c23-a004-3da25e5b8ce3', '0', '0'),
    ('d03307d4-2b28-4c23-a004-3da32e5b8a61', 'd03307d4-2b28-4c23-a004-3da25e5b8bb1', 'd03307d4-2b28-4c23-a004-3da25e5b8aa3', '0', '0')
    ON CONFLICT DO NOTHING;

INSERT INTO registros (registro_id, processo_id, placa, tipo_veiculo, cor_veiculo, marca_veiculo, armazenamento, confianca, criado_em) VALUES
    ('d03307d4-2b28-4d23-a003-2da32e5b8bc1', 'd03307d4-2b28-4c23-a004-3da32e5b8bb1', 'ABC2122', 'sedan', 'branco', 'honda', '1633298496_d03307d4-2b28-4d23-a003-2da32e5b8bb1', '0.5', NOW()),
    ('d03307d4-2b28-4d23-a003-3da32e5b8a61', 'd03307d4-2b28-4c23-a004-3da32e5b8bb1', 'ATX0220', 'sedan', 'preto', 'hyundai', '1633298496_d03307d4-2b28-4d23-a003-3da32e5b8a61', '0.5', NOW())
    ON CONFLICT DO NOTHING;

INSERT INTO veiculos (veiculo_id, placa, tipo, cor, marca, info) VALUES
    ('d03307d4-2b28-4d23-a004-3da32e5b8bb1', 'ABC0015', 'sedan', 'branco', 'honda', 'teste de carro 1'),
    ('d03307d4-2b28-4d23-a004-3da32e5b8a61', 'ABC0001', 'suv', 'preto', 'hyundai', 'teste de carro 2')
    ON CONFLICT DO NOTHING;
