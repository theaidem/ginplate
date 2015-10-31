console.log 'Go ahead for develop!'

Router			 	= ReactRouter
Link			 	= Router.Link
Route			 	= Router.Route
RouteHandler 		= Router.RouteHandler
NotFoundRoute	= Router.NotFoundRoute

App = React.createClass
	displayName: "App"


	componentDidMount:->

	onMem:(e)->

	render:->
		<div className="ui thesys container">
			<Menu/>
			<RouteHandler/>
		</div>


Menu = React.createClass
	render:()->
		<div className="ui stackable menu">
			<div className="header item">
				<i className="setting loading large green icon"></i>thesys
			</div>
			<Link to="about" className="item" >About</Link>
			<a className="item">Sign-in</a>
		</div>

About = React.createClass
	render:()->
		<div>About</div>

NotFound = React.createClass
	render:()->
		<div>NotFound</div>


routes = (
	<Route path="/" handler={App}>
		<Route path="about" name="about" handler={About} />
		<NotFoundRoute handler={NotFound} />
	</Route>	
)

Router.run(routes, Router.HistoryLocation,  (Handler)->
  	ReactDOM.render(<Handler/>, document.getElementById 'app')
)
