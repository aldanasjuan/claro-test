<script>
    import {url, session} from './data'
    let nit = ""

    let data
    let error = ""
    let loading = false
    async function verify(){
        if(nit == ""){
            error = "Ingresa un nit"
            return
        }
        loading = true
        let options = {
            method: "POST",
            headers: {
                "content-type": "application/json",
                "authorization": $session,
            },
            body:JSON.stringify({nit})
        }
        let res = await fetch(url+"/verify", options)
        if(res.status === 200){
            data = await res.json()
        }else{
            error = await res.text()
        }
        loading = false
    }
    let a = { "xml_name": { "Space": "", "Local": "return" }, "birthdate": "No disponible", "gender": "No disponible", "name": "No disponible", "nit": "1972656-2", "res": "correcto" }
</script>



<label>
    Verificar nit
    <input bind:value={nit}/>
</label>
<button disabled={loading} on:click={verify}>
    Verificar
</button>

<p style="color:red">{error}</p>
{#if data}
    <h3>Nombre</h3>
    <p>{data.name}</p>
    <h3>Fecha de nacimiento</h3>
    <p>{data.birthdate}</p>
    <h3>Género</h3>
    <p>{data.gender}</p>
    <h3>Nit</h3>
    <p>{data.nit}</p>
    <h3>Resultado</h3>
    <p>{data.res}</p>
{/if}