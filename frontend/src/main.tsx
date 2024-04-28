import React from 'react'
import ReactDOM from 'react-dom/client'
import './index.css'
import App from '@/app'
import StoreProvider from '@/store/store-provider'


ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <StoreProvider>
      <App/>
    </StoreProvider>
  </React.StrictMode>,
)
