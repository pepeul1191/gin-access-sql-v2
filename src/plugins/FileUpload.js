export class FileUpload {
  constructor(options = {}) {
    // Configurable DOM element IDs
    this.containerId = options.containerId || "fileUploadContainer";
    this.fileInputId = options.fileInputId || "imageFile";
    this.uploadButtonId = options.uploadButtonId || "uploadButton";
    this.viewButtonId = options.viewButtonId || "viewImageButton";
    this.errorMessageId = options.errorMessageId || "errorMessage";
    this.successMessageId = options.successMessageId || "successMessage";
    this.helpTextId = options.helpTextId || "helpText";
    this.labelElementId = options.labelElementId || null; // Optional
    
    // Component configuration
    this.label = options.label || "Seleccionar imagen";
    this.acceptedFormats = options.acceptedFormats || ["jpg", "jpeg", "png"];
    this.maxSizeMB = options.maxSizeMB || 2;
    this.baseURL = options.baseURL || "";
    this.url = this.baseURL + "/" + options.url || "/api/upload";
    this.fileKey = options.fileKey || "file";
    this.extraParams = options.extraParams || {};
    this.jwt = options.jwt || "";
    
    // Callbacks
    this.onSuccess = options.onSuccess || ((response) => console.log("Upload success:", response));
    this.onError = options.onError || ((error) => console.error("Upload error:", error));
    this.onViewClick = options.onViewClick || (() => {
      e.preventDefault();
    });
    
    this.file = null;
    this.isLoading = false;
    
    this.initElements();
    this.setupEventListeners();
    this.updateUI();
  }
  
  initElements() {
    // Get all DOM elements by configurable IDs
    this.container = document.getElementById(this.containerId);
    this.fileInput = document.getElementById(this.fileInputId);
    this.uploadButton = document.getElementById(this.uploadButtonId);
    this.viewButton = document.getElementById(this.viewButtonId);
    this.errorMessage = document.getElementById(this.errorMessageId);
    this.successMessage = document.getElementById(this.successMessageId);
    this.helpText = document.getElementById(this.helpTextId);
    
    // Get label element (either by ID or as first label in container)
    this.labelElement = this.labelElementId 
      ? document.getElementById(this.labelElementId)
      : this.container.querySelector('label');
    
    // Update initial values
    if (this.labelElement) {
      this.labelElement.textContent = this.label;
    }
    
    this.fileInput.accept = this.acceptedFormats.map(f => `.${f}`).join(',');
    this.helpText.textContent = 
      `Formatos aceptados: ${this.acceptedFormats.join(", ").toUpperCase()} (Máx. ${this.maxSizeMB}MB)`;
  }
  
  setupEventListeners() {
    this.fileInput.addEventListener('change', this.handleFileChange.bind(this));
    this.uploadButton.addEventListener('click', this.uploadFile.bind(this));
    this.viewButton.addEventListener('click', () => {
      /*if (this.file) {
        this.onViewClick();
      } else {
        this.showError("No hay imagen seleccionada para visualizar");
      }*/
      this.onViewClick();
    });
  }
  
  handleFileChange(event) {
    const selectedFile = event.target.files[0];
    if (!selectedFile) return;

    // Validate format
    const fileExt = selectedFile.name.split('.').pop().toLowerCase();
    if (!this.acceptedFormats.includes(fileExt)) {
      this.showError(`Formato no válido. Use: ${this.acceptedFormats.join(", ").toUpperCase()}`);
      return;
    }

    // Validate size
    if (selectedFile.size > this.maxSizeMB * 1024 * 1024) {
      this.showError(`El archivo excede el tamaño máximo de ${this.maxSizeMB}MB`);
      return;
    }

    this.clearMessages();
    this.file = selectedFile;
    this.updateUI();
  }
  
  async uploadFile() {
    if (!this.file) {
      this.showError("Por favor seleccione un archivo");
      return;
    }

    this.isLoading = true;
    this.clearMessages();
    this.updateUI();

    const formData = new FormData();
    formData.append(this.fileKey, this.file);

    // Add extra parameters
    Object.keys(this.extraParams).forEach(key => {
      formData.append(key, this.extraParams[key]);
    });

    try {
      const response = await fetch(this.url, {
        method: "POST",
        headers: {
          "Authorization": `Bearer ${this.jwt}`
        },
        body: formData, 
        mode: "cors", // Explicitamente habilitar CORS
        credentials: "include" // Solo si necesitas cookies
      });

      const data = await response.json();
      console.log(data)
      if (!response.ok) {
        throw new Error(data.message || "Error al subir el archivo");
      }

      this.showSuccess("Archivo subido correctamente");
      this.onSuccess(data);
    } catch (error) {
      this.showError(error.message);
      this.onError(error);
    } finally {
      this.isLoading = false;
      this.updateUI();
    }
  }
  
  showError(message) {
    this.errorMessage.textContent = message;
    this.successMessage.textContent = "";
  }
  
  showSuccess(message) {
    this.successMessage.textContent = message;
    this.errorMessage.textContent = "";
  }
  
  clearMessages() {
    this.errorMessage.textContent = "";
    this.successMessage.textContent = "";
  }
  
  updateUI() {
    // Update buttons state
    this.uploadButton.disabled = this.isLoading || !this.file;
    this.viewButton.disabled = this.isLoading || !this.file;
    
    // Update upload button text and icon
    const uploadIcon = this.uploadButton.querySelector('i');
    const uploadText = this.uploadButton.childNodes[2];
    
    if (this.isLoading) {
      uploadIcon.className = "fa fa-spinner spinner";
      uploadText.nodeValue = " Subiendo...";
    } else {
      uploadIcon.className = "fa fa-cloud-upload";
      uploadText.nodeValue = " Subir";
    }
    
    // Update input state
    this.fileInput.disabled = this.isLoading;
  }
}