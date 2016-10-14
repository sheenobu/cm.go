<!-- vim: syntax=html
-->
<addForm>
	<div id="addFormModal" class="modal">
		<div class="modal-content">
			<div class="row">
				<form class="col s12">
					<div class="row">
						<div class="input-field col s12">
							<input ref="name" id="name" type="text" class="validate"/>
							<label for="name">Name</label>
						</div>
					</div>
					<div class="row">
						<div class="input-field col s12">
							<input ref="description" id="description" type="text" class="validate"/>
							<label for="description">Description</label>
						</div>
					</div>
					<div class="row">
						<div class="input-field col s12">
							<input ref="url" id="url" type="text" class="validate"/>
							<label for="url">Project URL</label>
						</div>
					</div>
				</form>
			</div>
		</div>
		<div class="modal-footer">
			<a href="#!" onclick={click} class="modal-action modal-close waves-effect waves-green btn">Add Framework</a>
			<a href="#!" onclick={reset} class="left fmodal-action modal-close waves-effect waves-green btn">Cancel</a>
		</div>
	</div>

	<script>
		click() {

			var name = this.name.value.trim();
			var desc = this.description.value.trim();
			var url = this.url.value.trim();

			if (!name || !desc || !url) {
				return;
			}

			this.opts.add({name: name, description: desc, url: url})

			this.reset()
		}

		reset() {
			this.name.value = "";
			this.description.value = "";
			this.url.value = "";

			this.update();
		}

</addForm>
