package com.greencompass.feature.onboarding

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VerificationCodeScreen(
    onBack: () -> Unit,
    onVerify: () -> Unit,
    onSendAgain: () -> Unit
) {
    var code by remember { mutableStateOf(List(6) { "" }) }

    GreenCompassScaffold(
        title = "",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(horizontal = AppSpacing.lg),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Spacer(modifier = Modifier.height(AppSpacing.xxl))

            Text(
                text = "Verify your account",
                style = GreenCompassTypography.headlineLarge,
                color = GreenCompassColors.Charcoal,
                modifier = Modifier.padding(bottom = AppSpacing.md)
            )

            Text(
                text = "We sent a code to\namina.njeri@example.com",
                style = GreenCompassTypography.bodyLarge,
                color = GreenCompassColors.MutedText,
                textAlign = TextAlign.Center,
                modifier = Modifier.padding(bottom = AppSpacing.xxl)
            )

            Row(
                horizontalArrangement = Arrangement.spacedBy(AppSpacing.sm),
                modifier = Modifier.padding(bottom = AppSpacing.xxl)
            ) {
                repeat(6) { index ->
                    OutlinedTextField(
                        value = code[index],
                        onValueChange = { newValue ->
                            if (newValue.length <= 1 && newValue.all { it.isDigit() }) {
                                val newCode = code.toMutableList()
                                newCode[index] = newValue
                                code = newCode
                            }
                        },
                        modifier = Modifier
                            .width(48.dp)
                            .height(56.dp),
                        shape = RoundedCornerShape(12.dp),
                        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
                        textStyle = LocalTextStyle.current.copy(
                            textAlign = TextAlign.Center,
                            fontSize = GreenCompassTypography.headlineMedium.fontSize
                        )
                    )
                }
            }

            PrimaryButton(
                text = "Verify",
                onClick = onVerify,
                modifier = Modifier.padding(bottom = AppSpacing.xl)
            )

            Spacer(modifier = Modifier.weight(1f))

            Text(
                text = "Didn't receive a code?",
                style = GreenCompassTypography.bodyMedium,
                color = GreenCompassColors.MutedText,
                modifier = Modifier.padding(bottom = AppSpacing.xs)
            )

            TextLinkButton(
                text = "Send again",
                onClick = onSendAgain
            )

            Spacer(modifier = Modifier.height(AppSpacing.xxl))
        }
    }
}
