import logo from './logo.svg'
import styles from './App.module.css'

import Chart from './components/ForecastChart'

function App() {
    return (
        <div class={styles.App}>
            <header class={styles.header}></header>
            <Chart title="It's all in the wrist!" />
        </div>
    )
}

export default App
