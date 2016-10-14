<!-- vim: syntax=html
-->
<frameworkRow>
	<li class="collection-item avatar">
		<i class="material-icons circle">folder</i>
		<span class="title"><b>{opts.data.Name}</b></span>
		<p>{desc}</p>
		<a href={opts.data.URL}>{opts.data.URL}</a>
		<a class={classes} onclick={delete}>
			<i class="material-icons">delete</i>
		</a>
	</li>

	<script>
		this.classes = "btn-floating btn-small waves-effect secondary-content waves-light red";

		this.desc = this.opts.data.Description;
		if (this.desc.length > 100) {
			this.desc = this.desc.substr(0,100) + "..."
		}

		delete() {
			opts.delete(this.opts.data.ID)
		}
	</script>

</frameworkRow>
