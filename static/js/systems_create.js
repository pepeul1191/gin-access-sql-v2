document.addEventListener('DOMContentLoaded', function() {
  const today = new Date().toISOString().split('T')[0];
  const createdField = document.getElementById('created');
  if (createdField && !createdField.value) createdField.value = today;
  
  const updatedField = document.getElementById('updated');
  if (updatedField && !updatedField.value) updatedField.value = today;
});