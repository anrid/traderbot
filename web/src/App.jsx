import logo from './logo.svg'
import styles from './App.module.css'

import Hello from './components/Hello'

function App() {
    return (
        <div class={styles.App}>
            <header class={styles.header}>
                <img src={logo} class={styles.logo} alt="logo" />
                <p>
                    Edit <code>src/App.jsx</code> and save to reload.
                </p>
                <a
                    class={styles.link}
                    href="https://github.com/solidjs/solid"
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    Learn Solid
                </a>
                <p>
                    <Hello greeting="Dude!" />
                </p>
            </header>
        </div>
    )
}

export default App
