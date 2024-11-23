CREATE TYPE model_type AS ENUM ('device', 'service');
CREATE TYPE model_category AS ENUM ('wearable', 'camera', 'weather', 'entertainment');
CREATE TYPE protocol_type AS ENUM ('rest', 'grpc', 'mqtt', 'websocket');

CREATE TABLE smart_models (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type model_type NOT NULL,
    category model_category NOT NULL,
    manufacturer VARCHAR(255),
    model_number VARCHAR(100),
    metadata JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE smart_features (
    id UUID PRIMARY KEY,
    model_id UUID NOT NULL REFERENCES smart_models(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    protocol protocol_type NOT NULL,
    interface_path VARCHAR(255) NOT NULL,
    parameters JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_smart_models_type_category ON smart_models(type, category);
CREATE INDEX idx_smart_features_model_id ON smart_features(model_id);
