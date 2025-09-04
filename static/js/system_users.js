const alertDiv = document.getElementById('message');

const hideAlert = () =>{
  setTimeout(() => {
    alertDiv.classList.add('d-none');
  }, 5000);
}

document.addEventListener('DOMContentLoaded', function() {
  // Checkbox para seleccionar/deseleccionar todos
  const selectAllCheckbox = document.getElementById('select-all');
  const employeeCheckboxes = document.querySelectorAll('.user-checkbox');

  // Evento para el checkbox "Seleccionar todos"
  selectAllCheckbox.addEventListener('change', function() {
    employeeCheckboxes.forEach(checkbox => {
      checkbox.checked = selectAllCheckbox.checked;
    });
  });
  
  // Evento para los checkboxes individuales
  employeeCheckboxes.forEach(checkbox => {
    checkbox.addEventListener('change', function() {
      if (!this.checked) {
        selectAllCheckbox.checked = false;
      } else {
        // Verificar si todos los checkboxes están seleccionados
        const allChecked = Array.from(employeeCheckboxes).every(cb => cb.checked);
        selectAllCheckbox.checked = allChecked;
      }
    });
  });

  // Botón Guardar Cambios
  document.querySelector('.btn-success').addEventListener('click', function(e) {
    e.preventDefault();
    // Crear lista de employees con su estado de selección
    const employeesData = Array.from(employeeCheckboxes).map(checkbox => {
      return {
        id: parseInt(checkbox.value),
        selected: checkbox.checked
      };
    });


    // Enviar datos al servidor
    fetch(`/systems/${SYSTEM_ID}/users`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-CSRFToken': CSRF_TOKEN
      },
      body: JSON.stringify(employeesData)
    })
    .then(response => {
      alertDiv.classList.remove('d-none');
      if (response.ok) {
        //window.location.reload(); // Recargar la página después de guardar
        alertDiv.textContent = 'Cambios guardados correctamente';
        alertDiv.classList.add('alert-success');
      } else {
        alertDiv.textContent = 'Error al guardar los cambios';
        alertDiv.classList.add('alert-danger');
      }
      hideAlert();
    })
    .catch(error => {
      alertDiv.classList.remove('d-none');
      alertDiv.classList.add('alert-danger');
      alertDiv.textContent = 'Error al conectar con el servidor';
      console.error('Error:', error);
      hideAlert();
    });
  });
});