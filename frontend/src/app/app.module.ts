import { BrowserModule } from '@angular/platform-browser';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgModule } from '@angular/core';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { PostsComponent } from './posts/posts.component';
import { CreatePostComponent } from './createpost/createpost.component';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { ProfileComponent } from './profile/profile.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { httpInterceptorProviders } from './http-interceptors/index';
import { NoAuthGuard } from './no-auth.guard';
import { GatewayService } from './gateway.service';
import { HomeComponent } from './home/home.component';
import { GetPostComponent } from './get-post/get-post.component';
import { HeaderComponent } from './header/header.component';
import {MatToolbarModule} from '@angular/material/toolbar';
import {MatInputModule} from '@angular/material/input';
import {MatListModule} from '@angular/material/list';
import {MatIconModule} from '@angular/material/icon';
import {MatTableModule} from '@angular/material/table';
import {MatButtonModule} from '@angular/material/button';
import {MatMenuModule} from '@angular/material/menu';
import {MatPaginatorModule} from '@angular/material/paginator';
import { SearchComponent } from './search/search.component';
import {MatCardModule} from '@angular/material/card';
import { UserPostsComponent } from './user-posts/user-posts.component';
import { AuthGuard } from './auth.guard';
import { InterfacesService } from './interfaces.service';

@NgModule({
  declarations: [
    AppComponent,
    PostsComponent,
    CreatePostComponent,
    LoginComponent,
    RegisterComponent,
    ProfileComponent,
    HomeComponent,
    GetPostComponent,
    HeaderComponent,
    SearchComponent,
    UserPostsComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    AppRoutingModule,
    FormsModule,
    ReactiveFormsModule,
    BrowserAnimationsModule,
    MatToolbarModule,
    MatInputModule,
    MatListModule,
    MatIconModule,
    MatTableModule,
    MatButtonModule,
    MatMenuModule,
    MatPaginatorModule,
    MatCardModule,
  ],
  providers: [
    httpInterceptorProviders,
    NoAuthGuard,
    AuthGuard,
    GatewayService,
    InterfacesService,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }