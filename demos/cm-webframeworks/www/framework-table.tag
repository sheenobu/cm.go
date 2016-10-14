<!-- vim: syntax=html
-->
<frameworkTable>
	<ul class="collection">
		<li each={ item, i in opts.items }>
			<frameworkRow data={item} delete={parent.delete} />
		</li>
	</ul>

	<script>
		this.delete = this.opts.delete
	</script>

</frameworkTable>

