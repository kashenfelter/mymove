from locust import HttpLocust, TaskSet, task


class AnonBehavior(TaskSet):

    @task(1)
    def index(self):
        self.client.get("/")


class AnonUser(HttpLocust):
    task_set = AnonBehavior
    min_wait = 5000
    max_wait = 9000


class UserBehavior(TaskSet):

    csrf = None
    user = None
    token = None

    def _get_csrf_token(self):
        """
        Pull the CSRF token from the website by hitting the root URL.

        The token is set as a cookie with the name `masked_gorilla_csrf`
        """
        if self.csrf:
            return self.csrf
        self.client.get('/')
        self.csrf = self.client.cookies.get('masked_gorilla_csrf')

    def on_start(self):
        """ on_start is called when a Locust start before any task is scheduled """
        self._get_csrf_token()
        self.login()

    def on_stop(self):
        """ on_stop is called when the TaskSet is stopping """
        self.logout()

    def login(self):
        resp = self.client.post('/devlocal-auth/create',
                                headers={'x-csrf-token': self.csrf})
        try:
            self.user = resp.json()
            self.token = self.client.cookies.get('mil_session_token')
        except Exception:
            print('CSRF Token:', self.csrf)
            print(resp.content)

    def logout(self):
        self.client.post("/auth/logout")
        self.csrf = None
        self.user = None
        self.token = None

    @task(1)
    def index(self):
        self.client.get("/")


class MilMoveUser(HttpLocust):
    task_set = UserBehavior
    min_wait = 5000
    max_wait = 9000
