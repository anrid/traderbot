import styles from './Hello.module.css'

function Hello(props) {
    return <div class={styles.Hello}> Hello, {props.greeting}</div>
}

export default Hello
