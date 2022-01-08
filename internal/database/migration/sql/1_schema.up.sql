CREATE TABLE usuarios (
    usuario_id UUID,
    email TEXT UNIQUE,
    funcao TEXT[],
    senha TEXT,

    PRIMARY KEY (usuario_id)
);

CREATE TABLE cameras (
    camera_id UUID, -- TODO: POSTGRES tem DEFAULT gen_random_uuid()
    descricao TEXT NOT NULL,
    endereco_ip TEXT NOT NULL UNIQUE,
    porta INTEGER NOT NULL,
    canal INTEGER NOT NULL,
    usuario TEXT NOT NULL,
    senha TEXT NOT NULL,
    latitude TEXT NOT NULL,
    longitude TEXT NOT NULL,

    PRIMARY KEY (camera_id)
);

CREATE TABLE servidores_gravacao (
    servidor_gravacao_id UUID,
    endereco_ip TEXT NOT NULL,
    porta INTEGER NOT NULL,
    armazenamento TEXT NOT NULL,
    horas_retencao INTEGER NOT NULL,

    UNIQUE (endereco_ip, porta),
    PRIMARY KEY (servidor_gravacao_id)
);

CREATE TABLE processos (
    processo_id UUID,
    servidor_gravacao_id UUID,
    camera_id UUID,
    processador SMALLINT NOT NULL,
    adaptador SMALLINT NOT NULL,

    PRIMARY KEY (processo_id),
    FOREIGN KEY (servidor_gravacao_id) REFERENCES servidores_gravacao(servidor_gravacao_id) ON DELETE CASCADE
);

CREATE TABLE registros (
    registro_id UUID,
    processo_id UUID NOT NULL,
    placa TEXT NOT NULL,
    tipo_veiculo TEXT NOT NULL,
    cor_veiculo TEXT NOT NULL,
    marca_veiculo TEXT NOT NULL,
    armazenamento TEXT NOT NULL,
    confianca DECIMAL NOT NULL,
    criado_em TIMESTAMP WITH TIME ZONE NOT NULL,


    PRIMARY KEY (registro_id),
    FOREIGN KEY (processo_id) REFERENCES processos(processo_id) ON DELETE NO ACTION
);

CREATE TABLE veiculos (
    veiculo_id UUID,
    placa TEXT NOT NULL UNIQUE,
    tipo TEXT NOT NULL,
    cor TEXT NOT NULL,
    marca TEXT NOT NULL,
    info TEXT NOT NULL,

    UNIQUE (placa),
    PRIMARY KEY (veiculo_id)
);