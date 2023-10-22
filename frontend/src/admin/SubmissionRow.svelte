<script>
  import api from "../api";

  export let submission;
  let content = "";
  let dialog;
</script>

<tr>
  <td>{submission.username}</td>
  <td>{submission.title}</td>
  <td>{submission.timestamp}</td>
  <td>{submission.verdict}</td>
  <td
    ><button
      on:click={async () => {
        content = (await api.GET("/submissions/" + submission.ID, {})).body
          .content;
        dialog.showModal();
      }}>Show</button
    >
    <dialog on:close bind:this={dialog}>
      <pre><code>{content}</code></pre>
      <div>
        <button
          on:click={() => {
            dialog.close();
          }}>Close</button
        >
      </div>
    </dialog>
  </td></tr
>
