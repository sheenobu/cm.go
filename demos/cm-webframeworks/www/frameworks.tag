<!-- vim: syntax=html
-->
<frameworks>
    <div class="container" style="margin-top: 20px">
        <div class="row valign-wrapper">
            <div class="col s10 valign">
              <ul class="pagination left">
                  <li class={ disabled: canPrevious, waves-effect: !canPrevious }
					  onclick={onPrev}>
                    <a href="#prev"><i class="material-icons">chevron_left</i></a>
                </li>

                <li each={ n in nodes }
					data-page={n.i}
					onclick={parent.onPageChange}
                    class={ active: n.active, waves-effect: !n.active }>
                        {n.i+1}
                </li>

                <li class={ disabled: canNext, waves-effect: !canNext}
					onclick={onNext}>
                    <a href="#next"><i class="material-icons">chevron_right</i></a>
                </li>
              </ul>
            </div>
            <div class="col s2 valign">
                <a class="modal-trigger right waves-effect waves-light btn" href="#addFormModal">
                    <i class="material-icons left">add</i>Add
                </a>
            </div>
        </div>
        <div class="row">
            <div class="col s12">
                <frameworkTable
					items={ items }
					delete={ frameworkDelete }
				/>
            </div>
        </div>
		<addForm add={frameworkAdd} />
    </div>

    <script>

		this.currentPage = 0;

		onError(a, b, err) {
		    Materialize.toast(b, 4000) // 4000 is the duration of the toast
		}

		onMount() {
			$('.modal-trigger').leanModal();
			this.frameworkLoad();
		}

		onPageChange(evt) {
			this.currentPage = $(evt.toElement).data().page
			this.frameworkLoad()
		}

		onNext(evt) {
			this.currentPage = this.currentPage+1
			this.frameworkLoad()
		}

		onPrev(evt) {
			this.currentPage = this.currentPage-1
			this.frameworkLoad()
		}

        onData(data) {
            this.items = data.Frameworks;

			this.currentPage = data.CurrentPage
            this.canPrevious = data.CurrentPage != 1
            this.canNext = data.CurrentPage == data.PageCount-1

			this.nodes = [];
            for(var i = 0; i != data.PageCount; i++) {
                this.nodes[i] = {
                    i: i,
                    active: i == data.CurrentPage
                };
            }

            this.update();
        }

		frameworkDelete(frameworkID) {
			this.storage.remove(frameworkID, this.frameworkLoad)
		}

		frameworkAdd(framework) {
		    this.storage.add(framework, this.frameworkLoad)
		}

		frameworkLoad() {
		    this.storage.load(this.currentPage, this.onData);
		}

		this.storage = new Storage(this.onError);
		this.on("mount", this.onMount);

       </script>

</frameworks>
