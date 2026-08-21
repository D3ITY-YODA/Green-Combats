package com.greencompass.navigation

import androidx.compose.runtime.Composable
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import androidx.navigation.toRoute
import com.greencompass.feature.explore.*
import com.greencompass.feature.help.AboutScreen
import com.greencompass.feature.help.HelpScreen
import com.greencompass.feature.onboarding.*
import com.greencompass.feature.organizations.*
import com.greencompass.feature.places.*
import com.greencompass.feature.profile.*
import com.greencompass.feature.reports.*
import com.greencompass.feature.today.TodayRoute
import com.greencompass.feature.updates.*

@Composable
fun AppNavHost() {
    val navController = rememberNavController()

    NavHost(navController = navController, startDestination = AppRoute.Welcome) {
        // Onboarding
        composable<AppRoute.Welcome> { WelcomeScreen(onGetStarted = { navController.navigate(AppRoute.AccountChoice) }, onChooseLanguage = { navController.navigate(AppRoute.LanguageSelection) }) }
        composable<AppRoute.LanguageSelection> { LanguageSelectionScreen(onContinue = { navController.navigate(AppRoute.AccountChoice) }, onBack = { navController.popBackStack() }) }
        composable<AppRoute.AccountChoice> { AccountChoiceScreen(onPersonal = { navController.navigate(AppRoute.PersonalRegistration) }, onOrganization = { navController.navigate(AppRoute.OrganizationSearch) }, onInvitation = { navController.navigate(AppRoute.InvitationAcceptance) }, onSignIn = { navController.navigate(AppRoute.PhoneEmailSignIn) }) }
        composable<AppRoute.PersonalRegistration> { PersonalRegistrationScreen(onBack = { navController.popBackStack() }, onContinueWithGoogle = { navController.navigate(AppRoute.GoogleSignIn) }, onContinue = { navController.navigate(AppRoute.VerificationCode) }) }
        composable<AppRoute.GoogleSignIn> { GoogleSignInScreen(onBack = { navController.popBackStack() }, onSuccess = { navController.navigate(AppRoute.VerificationCode) }, onUsePhoneEmail = { navController.popBackStack(); navController.navigate(AppRoute.PhoneEmailSignIn) }) }
        composable<AppRoute.PhoneEmailSignIn> { PhoneEmailSignInScreen(onBack = { navController.popBackStack() }, onContinue = { navController.navigate(AppRoute.VerificationCode) }, onCreateAccount = { navController.popBackStack(); navController.navigate(AppRoute.PersonalRegistration) }) }
        composable<AppRoute.VerificationCode> { VerificationCodeScreen(onBack = { navController.popBackStack() }, onVerify = { navController.navigate(AppRoute.LocationSelection) }, onSendAgain = { }) }
        composable<AppRoute.OrganizationSearch> { OrganizationSearchScreen(onBack = { navController.popBackStack() }, onSelectOrganization = { navController.navigate(AppRoute.OrganizationSelection) }, onRequestAccess = { navController.navigate(AppRoute.RequestAccess) }) }
        composable<AppRoute.OrganizationSelection> { OrganizationSelectionScreen(onBack = { navController.popBackStack() }, onContinue = { navController.navigate(AppRoute.AccessPending) }) }
        composable<AppRoute.RequestAccess> { RequestAccessScreen(onBack = { navController.popBackStack() }, onRequestSent = { navController.navigate(AppRoute.AccessPending) }) }
        composable<AppRoute.InvitationAcceptance> { InvitationAcceptanceScreen(onBack = { navController.popBackStack() }, onAccept = { navController.navigate(AppRoute.LocationSelection) }, onDecline = { navController.popBackStack() }) }
        composable<AppRoute.AccessPending> { AccessPendingScreen(onBack = { navController.popBackStack() }, onViewRequest = { }, onReturnToToday = { navController.navigate(AppRoute.Today) }) }
        composable<AppRoute.LocationSelection> { LocationSelectionScreen(onBack = { navController.popBackStack() }, onUseLocation = { navController.navigate(AppRoute.Interests) }, onSearchPlace = { navController.navigate(AppRoute.SearchPlace) }, onChooseOnMap = { navController.navigate(AppRoute.MapPlaceSelection) }, onContinue = { navController.navigate(AppRoute.Interests) }) }
        composable<AppRoute.SearchPlace> { SearchPlaceScreen(onBack = { navController.popBackStack() }, onSelectPlace = { navController.popBackStack(); navController.navigate(AppRoute.Interests) }) }
        composable<AppRoute.MapPlaceSelection> { MapPlaceSelectionScreen(onBack = { navController.popBackStack() }, onConfirmPlace = { navController.popBackStack(); navController.navigate(AppRoute.Interests) }) }
        composable<AppRoute.Interests> { InterestsScreen(onBack = { navController.popBackStack() }, onContinue = { navController.navigate(AppRoute.NotificationPreferences) }, onSkip = { navController.navigate(AppRoute.NotificationPreferences) }) }
        composable<AppRoute.NotificationPreferences> { NotificationPreferencesScreen(onBack = { navController.popBackStack() }, onContinue = { navController.navigate(AppRoute.PrivacyPermission) }) }
        composable<AppRoute.PrivacyPermission> { PrivacyPermissionScreen(onBack = { navController.popBackStack() }, onAllowLocation = { navController.navigate(AppRoute.SetupComplete) }, onChooseManually = { navController.navigate(AppRoute.SetupComplete) }, onReadPrivacy = { }) }
        composable<AppRoute.SetupComplete> { SetupCompleteScreen(onFinish = { navController.navigate(AppRoute.Today) { popUpTo(AppRoute.Welcome) { inclusive = true } } }) }

        // Main App
        composable<AppRoute.Today> { TodayRoute() }
        composable<AppRoute.Explore> { ExploreScreen(onNavigate = { route -> navController.navigate(route) }) }
        composable<AppRoute.LocalOutlook> { LocalOutlookScreen(onBack = { navController.popBackStack() }) }
        composable<AppRoute.SeasonalInformation> { SeasonalInformationScreen(onBack = { navController.popBackStack() }) }
        composable<AppRoute.WaterOutlook> { WaterOutlookScreen(onBack = { navController.popBackStack() }) }
        composable<AppRoute.LandEcosystems> { LandEcosystemsScreen(onBack = { navController.popBackStack() }) }
        composable<AppRoute.FoodAgriculture> { FoodAgricultureScreen(onBack = { navController.popBackStack() }) }
        composable<AppRoute.CommunityUpdates> { CommunityUpdatesScreen(onBack = { navController.popBackStack() }) }

        // Updates & Reports
        composable<AppRoute.Updates> { UpdatesScreen(onNavigate = { route -> navController.navigate(route) }) }
        composable<AppRoute.UpdateDetail> { backStackEntry ->
            val route = backStackEntry.toRoute<AppRoute.UpdateDetail>()
            UpdateDetailScreen(updateId = route.updateId, onBack = { navController.popBackStack() }, onAcknowledge = { navController.navigate(AppRoute.UpdateAcknowledgement) }, onShare = { navController.navigate(AppRoute.Report) })
        }
        composable<AppRoute.UpdateAcknowledgement> { UpdateAcknowledgementScreen(onBack = { navController.popBackStack() }, onShare = { navController.navigate(AppRoute.Report) }) }
        composable<AppRoute.Report> { ReportScreen(onSelectType = { type -> navController.navigate(AppRoute.ReportTypeSelection(type)) }) }
        composable<AppRoute.ReportTypeSelection> { backStackEntry ->
            val route = backStackEntry.toRoute<AppRoute.ReportTypeSelection>()
            ReportTypeSelectionScreen(reportType = route.reportType, onBack = { navController.popBackStack() }, onContinue = { navController.navigate(AppRoute.ReportForm(route.reportType)) }, onChooseDifferent = { navController.popBackStack() })
        }
        composable<AppRoute.ReportForm> { backStackEntry ->
            val route = backStackEntry.toRoute<AppRoute.ReportForm>()
            ReportFormScreen(reportType = route.reportType, onBack = { navController.popBackStack() }, onSubmit = { navController.navigate(AppRoute.ReportSubmitted) })
        }
        composable<AppRoute.ReportSubmitted> { ReportSubmittedScreen(onViewStatus = { navController.navigate(AppRoute.ReportStatus) }, onReturn = { navController.navigate(AppRoute.Today) { popUpTo(AppRoute.Today) { inclusive = true } } }) }
        composable<AppRoute.ReportStatus> { ReportStatusScreen(onBack = { navController.popBackStack() }) }

        // Places
        composable<AppRoute.Places> { PlacesScreen(onBack = { navController.popBackStack() }, onAddPlace = { navController.navigate(AppRoute.AddPlace) }, onOpenPlace = { placeId -> navController.navigate(AppRoute.PlaceDetails(placeId)) }) }
        composable<AppRoute.AddPlace> { AddPlaceScreen(onBack = { navController.popBackStack() }, onUseLocation = { }, onSearch = { }, onChooseOnMap = { }) }
        composable<AppRoute.PlaceDetails> { backStackEntry ->
            val route = backStackEntry.toRoute<AppRoute.PlaceDetails>()
            PlaceDetailsScreen(placeId = route.placeId, onBack = { navController.popBackStack() }, onSetPrimary = { }, onEditName = { }, onRemove = { })
        }

        // Profile & Settings
        composable<AppRoute.Profile> { ProfileScreen(onNavigate = { route -> navController.navigate(route) }) }
        composable<AppRoute.Settings> { SettingsScreen(onBack = { navController.popBackStack() }, onNavigate = { route -> navController.navigate(route) }, onSignOut = { }) }
        composable<AppRoute.NotificationSettings> { NotificationSettingsScreen(onBack = { navController.popBackStack() }) }
        composable<AppRoute.LanguageSettings> { LanguageSettingsScreen(onBack = { navController.popBackStack() }) }
        composable<AppRoute.AccessibilitySettings> { AccessibilitySettingsScreen(onBack = { navController.popBackStack() }) }
        composable<AppRoute.PrivacySettings> { PrivacySettingsScreen(onBack = { navController.popBackStack() }, onDeleteAccount = { }) }

        // Organizations
        composable<AppRoute.Organizations> { OrganizationsScreen(onBack = { navController.popBackStack() }, onOpenOrganization = { navController.navigate(AppRoute.OrganizationView) }, onViewRequest = { }) }
        composable<AppRoute.OrganizationSwitcher> { OrganizationSwitcherScreen(onBack = { navController.popBackStack() }, onSelect = { navController.popBackStack() }) }
        composable<AppRoute.OrganizationView> { OrganizationViewScreen(onBack = { navController.popBackStack() }, onOpenOnWeb = { }, onReturnToPersonal = { navController.popBackStack() }) }

        // Help & About
        composable<AppRoute.Help> { HelpScreen(onBack = { navController.popBackStack() }, onContactSupport = { }) }
        composable<AppRoute.About> { AboutScreen(onBack = { navController.popBackStack() }) }
    }
}
