import logo from './logo.svg'
import styles from './App.module.css'

import Chart from './components/LineRacerChart'

function App() {
    return (
        <div class={styles.App}>
            <header class={styles.header}></header>
            <Chart title="Income dude!" />
        </div>
    )
}

export default App
