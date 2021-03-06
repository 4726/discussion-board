import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { CreatePostComponent } from './createpost/createpost.component';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { NoAuthGuard } from './no-auth.guard';
import { AuthGuard } from './auth.guard';
import { HomeComponent } from './home/home.component';
import { PostsComponent } from './posts/posts.component';
import { GetPostComponent } from './get-post/get-post.component';
import { SearchComponent } from './search/search.component';
import { ProfileComponent } from './profile/profile.component';
import { UserPostsComponent } from './user-posts/user-posts.component';

const routes: Routes = [
  {path: 'post/create', component: CreatePostComponent, canActivate: [AuthGuard]},
  {path: 'login', component: LoginComponent, canActivate: [NoAuthGuard]},
  {path: 'register', component: RegisterComponent, canActivate: [NoAuthGuard]},
  {path: 'home', component: HomeComponent},
  {path: 'posts/:page', component: PostsComponent},
  {path: 'post/:postID', component: GetPostComponent},
  {path: 'search', component: SearchComponent},
  {path: 'profile/:userid', component: ProfileComponent},
  {path: 'profile/:userid/posts/:page', component: UserPostsComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
