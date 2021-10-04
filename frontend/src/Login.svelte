<script>
    import {url, session} from './data'
    let auth = {
        "username": "",
        "password": ""
    }
    let isRegister = false
    let reg = {
        "username": "",
        "password": "",
        "firstName": "",
        "lastName": "",
        "location": "GT",
        "role": "user",
    }
    let loading = false
    let error = ""

    async function login(){
        loading = true
        let options = {
            method: "POST",
            headers: {
                "content-type": "application/json"
            },
            body:JSON.stringify(auth)
        }
        let res = await fetch(url+"/login", options)
        if(res.status === 200){
            let token = res.headers.get("authorization")
            if(token){
                session.set(token)
            }
        }else{
            error = await res.text()
        }
        loading = false
    }
    async function register(){
        loading = true
        let options = {
            method: "POST",
            headers: {
                "content-type": "application/json"
            },
            body:JSON.stringify(reg)
        }
        let res = await fetch(url+"/register", options)
        if(res.status === 200){
            let token = res.headers.get("authorization")
            if(token){
                session.set(token)
            }
        }else{
            error = await res.text()
        }
        loading = false
    }
</script>

{#if isRegister}
    <h2>
        Log in
    </h2>
    <label>
        First name
        <input bind:value={reg.firstName}/>
    </label>
    <label>
        Last Name
        <input bind:value={reg.lastName}/>
    </label>
    <label>
        Username
        <input bind:value={reg.username}/>
    </label>
    <label>
        Password
        <input bind:value={reg.password}/>
    </label>
    <button disabled={loading} on:click={register}>
        Register
    </button>

    <button on:click={() => isRegister = false}>
        Log in instead
    </button>

    <p style="color:red">
        {error}
    </p>
{:else}
    <h2>
        Log in
    </h2>
    <label>
        Username
        <input bind:value={auth.username}/>
    </label>
    <label>
        Password
        <input bind:value={auth.password}/>
    </label>
    <button disabled={loading} on:click={login}>
        Log in
    </button>
    <button on:click={() => isRegister = true}>
        Register instead
    </button>
    <p style="color:red">
        {error}
    </p>
{/if}