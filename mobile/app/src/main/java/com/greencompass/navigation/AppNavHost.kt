package com.greencompass.navigation

import androidx.compose.runtime.Composable
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.greencompass.feature.onboarding.*

@Composable
fun AppNavHost() {
    val navController = rememberNavController()

    NavHost(navController = navController, startDestination = AppRoute.Welcome) {
        composable<AppRoute.Welcome> {
            WelcomeScreen(
                onGetStarted = { navController.navigate(AppRoute.AccountChoice) },
                onChooseLanguage = { navController.navigate(AppRoute.LanguageSelection) }
            )
        }
        composable<AppRoute.LanguageSelection> {
            LanguageSelectionScreen(
                onContinue = { navController.navigate(AppRoute.AccountChoice) },
                onBack = { navController.popBackStack() }
            )
        }
        composable<AppRoute.AccountChoice> {
            AccountChoiceScreen(
                onPersonal = { navController.navigate(AppRoute.PersonalRegistration) },
                onOrganization = { /* TODO: Organization Search */ },
                onInvitation = { /* TODO: Invitation Acceptance */ },
                onSignIn = { navController.navigate(AppRoute.PhoneEmailSignIn) }
            )
        }
        composable<AppRoute.PersonalRegistration> {
            PersonalRegistrationScreen(
                onBack = { navController.popBackStack() },
                onContinueWithGoogle = { navController.navigate(AppRoute.GoogleSignIn) },
                onContinue = { navController.navigate(AppRoute.VerificationCode) } // Simplified flow for demo
            )
        }
        composable<AppRoute.GoogleSignIn> {
            GoogleSignInScreen(
                onBack = { navController.popBackStack() },
                onSuccess = { navController.navigate(AppRoute.VerificationCode) },
                onUsePhoneEmail = { 
                    navController.popBackStack() 
                    navController.navigate(AppRoute.PhoneEmailSignIn) 
                }
            )
        }
        composable<AppRoute.PhoneEmailSignIn> {
            PhoneEmailSignInScreen(
                onBack = { navController.popBackStack() },
                onContinue = { navController.navigate(AppRoute.VerificationCode) },
                onCreateAccount = { 
                    navController.popBackStack() 
                    navController.navigate(AppRoute.PersonalRegistration) 
                }
            )
        }
        composable<AppRoute.VerificationCode> {
            VerificationCodeScreen(
                onBack = { navController.popBackStack() },
                onVerify = { /* TODO: Navigate to Location Selection */ },
                onSendAgain = { /* TODO: Resend logic */ }
            )
        }
    }
}
