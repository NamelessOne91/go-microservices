// Create "logs" collection and insert test documents
db.createCollection("logs");

const logEntries = [
  {
    _id: "61799e857fb7ed9569be5ccf",
    name: 'Log Entry 1',
    data: 'Data for Log Entry 1',
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    _id: "61799e857fb7ed9569be5cd0",
    name: 'Log Entry 2',
    data: 'Data for Log Entry 2',
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    _id: "61799e857fb7ed9569be5cd1",
    name: 'Log Entry 3',
    data: 'Data for Log Entry 3',
    created_at: new Date(),
    updated_at: new Date()
  }
];

db.logs.insertMany(logEntries);