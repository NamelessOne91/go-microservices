// Create "logs" collection and insert test documents
db.createCollection("logs");

var logs = [
  {
    name: "Log 1",
    data: "Data 1",
    created_at: new ISODate(),
    updated_at: new ISODate()
  },
  {
    name: "Log 2",
    data: "Data 2",
    created_at: new ISODate(),
    updated_at: new ISODate()
  },
  {
    name: "Log 3",
    data: "Data 3",
    created_at: new ISODate(),
    updated_at: new ISODate()
  }
];

db.logs.insertMany(logs);