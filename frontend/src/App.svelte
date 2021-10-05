<script>
	import Nit from './Nit.svelte'
	import Login from './Login.svelte'
	import {session, user, url} from './data.js'
	console.log($session)

	let number = 0
	let error = ""
	let result = ""
	let loading = false
	async function logout(){
		let options = {
            headers: {
                "authorization": $session
            }
        }
        let res = await fetch(url+"/logout", options)
        if(res.status === 200){
            session.set(null)
        }else{
            error = await res.text()
        }
	}
	async function roman(){
		
		error = ""
		result = ""
		number = parseInt(number)
		if(number < 1){
			error = "No. debe ser mayor a 0"
			return
		}
		loading = true
		let options = {
			method: "POST",
            headers: {
				"content-type": "application/json"
            },
            body:JSON.stringify({number})
        }
		let res = await fetch(url+"/roman", options)
        if(res.status === 200){
			result = await res.text()
			
        }else{
			error = await res.text()
        }
		loading = false
	}
</script>



{#if $session}
	{#if $user && $user.firstName}
		<p>Hola, {$user.firstName}</p>
	{/if}
	<button on:click={logout}>Logout</button>
	
	<label>Número a Romano
		<input type="number" bind:value={number}/>
	</label>
	<button disabled={loading} on:click={roman}>Convertir</button>
	<p style="color:red">
		{error}
	</p>
	<p>
		{result}
	</p>
	<Nit/>
{:else}
	<Login/>
{/if}



