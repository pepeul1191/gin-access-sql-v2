// src/entries/vendor.js
// libs
import 'bootstrap/dist/css/bootstrap.min.css';
import 'font-awesome/css/font-awesome.min.css';
import * as bootstrap from 'bootstrap/dist/js/bootstrap.bundle.min.js';
import axios from 'axios'; // Simplified import

// Explicitly expose axios and bootstrap globally
window.axios = axios;
window.bootstrap = bootstrap;