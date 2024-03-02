CREATE TABLE logs(  
    id INT NOT NULL,
    project_id INT NOT NULL,
    name String NOT NULL,
    description String,
    priority INT NOT NULL,
    removed BOOLEAN NOT NULL,
    created_at DATETIME DEFAULT now()
)
ENGINE = MergeTree()
ORDER BY (id)
PRIMARY KEY(id);