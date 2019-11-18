import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { CreatePostComponent } from './createpost/createpost.component';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { NoAuthGuard } from './no-auth.guard';


const routes: Routes = [
  {path: 'post/create', component: CreatePostComponent},
  {path: 'login', component: LoginComponent, canActivate: [NoAuthGuard]},
  {path: 'register', component: RegisterComponent, canActivate: [NoAuthGuard]}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
