{{define "pagetitle"}}Root{{end}}

{{$loggedin := .loggedin}}
{{if $loggedin}}
<div class="row" style="margin-bottom: 20px;">
	<div class="col-md-offset-9 col-md-2 text-right">
		<a class="btn btn-primary" href="/blogs/new"><i class="fa fa-plus"></i> New Post</a>	
	</div>
</div>
{{end}}

<div class="row">
	<div class="col-md-offset-1 col-md-10">
		<div class="panel panel-info">
			<div class="panel-heading">
				Title
			</div>
			<div class="panel-body">Content</div>
			<div class="panel-footer">
				Footer
			</div>
		</div>
	</div>
</div>