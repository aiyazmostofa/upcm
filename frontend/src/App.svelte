<script>
  import Login from "./Login.svelte";
  import Team from "./team/Team.svelte";
  import Admin from "./admin/Admin.svelte";
  import stores from "./stores";

  let user = {
    token: "",
    authLevel: "",
    ID: "",
  };

  stores.user.subscribe((v) => {
    user = v;
  });
</script>

<main>
  {#if user.token}
    {#if user.authLevel == "Admin"}
      <Admin />
    {:else}
      <Team />
    {/if}
  {:else}
    <Login />
  {/if}
</main>
