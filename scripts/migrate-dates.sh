#!/bin/bash

# MongoDB Migration Script
# Converts string deadline fields to DateTime format in tasks collection

MONGO_URI="${MONGO_URI:-mongodb://localhost:27017}"
MONGO_DATABASE="${MONGO_DATABASE:-journaling}"

echo "🔄 Starting MongoDB migration..."
echo "📦 Database: $MONGO_DATABASE"
echo ""

# MongoDB shell command to convert string dates to DateTime
mongosh "$MONGO_URI/$MONGO_DATABASE" <<'EOF'

// Function to convert date string to ISODate
function convertStringToDate(dateStr) {
  if (!dateStr || typeof dateStr !== 'string') {
    return null;
  }

  try {
    // Parse the date string and create ISODate
    return new Date(dateStr);
  } catch (e) {
    print('Error parsing date: ' + dateStr);
    return null;
  }
}

// Find all tasks with string deadline field
print('🔍 Finding tasks with string deadline...');
const tasksWithStringDeadline = db.tasks.find({
  deadline: { $type: "string" }
}).toArray();

print('Found ' + tasksWithStringDeadline.length + ' tasks to migrate');

// Update each task
let updated = 0;
let failed = 0;

tasksWithStringDeadline.forEach(task => {
  const dateValue = convertStringToDate(task.deadline);

  if (dateValue) {
    try {
      db.tasks.updateOne(
        { _id: task._id },
        {
          $set: {
            deadline: dateValue,
            updated_at: new Date()
          }
        }
      );
      updated++;
      print('✅ Updated task: ' + task._id + ' (' + task.title + ')');
    } catch (e) {
      failed++;
      print('❌ Failed to update task: ' + task._id + ' - ' + e.message);
    }
  } else {
    failed++;
    print('❌ Invalid date format for task: ' + task._id);
  }
});

print('');
print('📊 Migration Summary:');
print('   ✅ Successfully updated: ' + updated);
print('   ❌ Failed: ' + failed);
print('   📝 Total: ' + tasksWithStringDeadline.length);

// Do the same for todos collection
print('');
print('🔍 Finding todos with string due_date...');
const todosWithStringDueDate = db.todos.find({
  due_date: { $type: "string" }
}).toArray();

print('Found ' + todosWithStringDueDate.length + ' todos to migrate');

let updatedTodos = 0;
let failedTodos = 0;

todosWithStringDueDate.forEach(todo => {
  const dateValue = convertStringToDate(todo.due_date);

  if (dateValue) {
    try {
      db.todos.updateOne(
        { _id: todo._id },
        {
          $set: {
            due_date: dateValue,
            updated_at: new Date()
          }
        }
      );
      updatedTodos++;
      print('✅ Updated todo: ' + todo._id + ' (' + todo.title + ')');
    } catch (e) {
      failedTodos++;
      print('❌ Failed to update todo: ' + todo._id + ' - ' + e.message);
    }
  } else {
    failedTodos++;
    print('❌ Invalid date format for todo: ' + todo._id);
  }
});

print('');
print('📊 Todos Migration Summary:');
print('   ✅ Successfully updated: ' + updatedTodos);
print('   ❌ Failed: ' + failedTodos);
print('   📝 Total: ' + todosWithStringDueDate.length);

print('');
print('✨ Migration completed!');

EOF

echo ""
echo "✅ Migration script completed"
echo ""
echo "To verify, run:"
echo "  mongosh $MONGO_URI/$MONGO_DATABASE --eval 'db.tasks.find({deadline: {\$exists: true}}).pretty()'"
