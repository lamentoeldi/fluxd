package wfsqlite

const CreateWorkflowsTable = `CREATE TABLE IF NOT EXISTS workflows (
    id TEXT NOT NULL PRIMARY KEY,
    content TEXT NOT NULL,
);`
